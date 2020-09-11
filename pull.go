package main

import (
	"context"
	"fmt"
	"github.com/cirocosta/grafana-sync/grafana"
	"github.com/pkg/errors"
	"os"
	"path"
)

type pullCommand struct{
	SyncFolders []string `long:"sync-folders" short:"f" description:"datasource used by the dashboards"`
}

func (p *pullCommand) Execute(args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)
	client := grafana.NewClient(grafana.ClientConfig{
		Address:     config.Address,
		Verbose:     config.Verbose,
		AccessToken: config.Auth.AccessToken,
		Username:    config.Auth.Username,
		Password:    config.Auth.Password,
	})

	refs, err := client.ListDashboardSearchEntries(ctx)
	if err != nil {
		return
	}

	var dashboard map[string]interface{}
	for _, ref := range refs {
		if ! shouldSyncFolder(p.SyncFolders, ref.Folder) {
			fmt.Println("Not syncing folder", ref.Folder)
			continue
		}
		dashboardFolderInFs := path.Join(string(config.Directory), ref.Folder)

		err = eventuallyCreateDirectory(dashboardFolderInFs)
		if err != nil {
			return
		}

		dashboard, err = client.GetDashboard(ctx, ref.Uid)
		if err != nil {
			return
		}

		err = grafana.SaveToDisk(path.Join(dashboardFolderInFs, ref.Title)+".json", dashboard)
		if err != nil {
			return
		}
	}

	return
}
func shouldSyncFolder(syncFolders []string, folder string) bool {
	for _, f := range(syncFolders) {
		if folder == f {
			return true
		}
	}
	return false
}

func eventuallyCreateDirectory(dir string) (err error) {
	finfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				err = errors.Wrapf(err,
					"failed to create directory %s", dir)
				return
			}

			return
		}

		err = errors.Wrapf(err,
			"unexpected failure checking directory %s", dir)
		return
	}

	if !finfo.IsDir() {
		err = errors.Errorf("location %s already exists and is not a directory", dir)
		return
	}

	return
}
