package main

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/cirocosta/grafana-sync/grafana"
	"github.com/pkg/errors"
)

type pushCommand struct {
	Datasource string `long:"data-source" short:"d" description:"datasource used by the dashboards"`
}

// for each directory:
//   - try to gather a folder ID for it
//	- if it doesn't exist: create a folder
//   - create the dashboards in the directory
func (p *pushCommand) Execute(args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	client := grafana.NewClient(grafana.ClientConfig{
		Address:     config.Address,
		Verbose:     config.Verbose,
		AccessToken: config.Auth.AccessToken,
		Username:    config.Auth.Username,
		Password:    config.Auth.Password,
	})

	entries, err := ioutil.ReadDir(string(config.Directory))
	if err != nil {
		return
	}

	var (
		folder           grafana.Folder
		folderEntries    []os.FileInfo
		dashboardEntries []grafana.DashboardCreateOrUpdateEntry
		existingFolders  = map[string]grafana.Folder{}
	)

	grafanaFolders, err := client.ListFolders(ctx)
	if err != nil {
		return
	}

	for _, grafanaFolder := range grafanaFolders {
		existingFolders[grafanaFolder.Title] = grafanaFolder
	}

	for _, entry := range entries {
		var (
			dashboard map[string]interface{}
			filepath  = path.Join(string(config.Directory), entry.Name())
		)

		if !entry.IsDir() {
			dashboard, err = grafana.LoadFromDisk(filepath)
			if err != nil {
				return
			}

			dashboardEntries = append(dashboardEntries, grafana.DashboardCreateOrUpdateEntry{
				Overwrite: true,
				FolderId:  0,
				Dashboard: dashboard,
			})

			continue
		}

		var folderAlreadyExists bool
		folder, folderAlreadyExists = existingFolders[entry.Name()]
		if !folderAlreadyExists {
			folder, err = client.CreateFolder(ctx, entry.Name())
			if err != nil {
				return
			}
		}

		folderEntries, err = ioutil.ReadDir(filepath)
		if err != nil {
			return
		}

		for _, folderEntry := range folderEntries {
			if folderEntry.IsDir() {
				err = errors.Errorf("can't have a dir within a grafana folder %s",
					folderEntry.Name())
				return
			}

			dashboard, err = grafana.LoadFromDisk(path.Join(filepath, folderEntry.Name()))
			if err != nil {
				return
			}

			dashboardEntries = append(dashboardEntries, grafana.DashboardCreateOrUpdateEntry{
				Overwrite: true,
				FolderId:  folder.Id,
				Dashboard: dashboard,
			})
		}
	}

	for _, dashboardEntry := range dashboardEntries {
		if p.Datasource != "" {
			err = grafana.SetPanelDatasources(
				dashboardEntry.Dashboard, p.Datasource)
			if err != nil {
				return
			}
		}

		err = client.PushDashboard(ctx, dashboardEntry)
		if err != nil {
			return
		}
	}

	return
}
