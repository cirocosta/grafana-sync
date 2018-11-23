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
		folder        grafana.Folder
		folderEntries []os.FileInfo
	)
	for _, entry := range entries {
		if !entry.IsDir() {
			err = pushDashboardFromDisk(ctx, client, path.Join(string(config.Directory), entry.Name()), 0)
			if err != nil {
				return
			}

			continue
		}

		folder, err = client.CreateFolder(ctx, entry.Name())
		if err != nil {
			return
		}

		folderEntries, err = ioutil.ReadDir(path.Join(string(config.Directory), entry.Name()))
		if err != nil {
			return
		}

		for _, folderEntry := range folderEntries {
			if entry.IsDir() {
				err = errors.Errorf("can't have a dir within a grafana folder")
				return
			}

			err = pushDashboardFromDisk(ctx, client, path.Join(
				string(config.Directory), entry.Name(), folderEntry.Name()), folder.Id)
			if err != nil {
				return
			}
		}
	}

	return
}

func pushDashboardFromDisk(ctx context.Context, c *grafana.Client, filepath string, folderId int) (err error) {
	dashboard, err := grafana.LoadFromDisk(filepath)
	if err != nil {
		err = errors.Wrapf(err,
			"could'nt load into mem")
		return
	}

	entry := grafana.DashboardCreateOrUpdateEntry{
		Overwrite: true,
		FolderId:  folderId,
		Dashboard: dashboard,
	}

	err = c.PushDashboard(ctx, entry)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to push dashboard")
		return
	}

	return
}
