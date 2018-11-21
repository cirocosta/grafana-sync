package main

import (
	"context"
	"os"
	"path"

	"github.com/cirocosta/grafana-sync/grafana"
	"github.com/pkg/errors"
)

type pullCommand struct{}

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

	refs, err := client.ListDashboardRefs(ctx)
	if err != nil {
		return
	}

	var dashboard grafana.Dashboard
	for _, ref := range refs {
		dashboardFolderInFs := path.Join(string(config.Directory), ref.Folder)

		err = eventuallyCreateDirectory(dashboardFolderInFs)
		if err != nil {
			return
		}

		dashboard, err = client.GetDashboard(ctx, ref.Uid)
		if err != nil {
			return
		}

		err = dashboard.SaveToDisk(path.Join(dashboardFolderInFs, ref.Title) + ".json")
		if err != nil {
			return
		}
	}

	return
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
