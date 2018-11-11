package main

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"path"

	"github.com/cirocosta/grafana-sync/grafana"
	"github.com/jessevdk/go-flags"
)

var config struct {
	AccessToken string        `long:"access-token" required:"true" description:"access token to authenticate against grafana"`
	Address     string        `long:"address" default:"http://localhost:3000" description:"grafana address"`
	Directory   DirectoryFlag `long:"directory" default:"./" description:"directory where dashboards live"`
	Verbose     bool          `short:"v" long:"verbose" description:"displays requests on stderr"`
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

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		return
	}

	client := grafana.NewClient(config.Address, config.AccessToken, config.Verbose)
	dashboards, err := client.AllDashboards(context.Background())
	if err != nil {
		panic(err)
	}

	for _, dashboard := range dashboards {
		dashboardFolderInFs := path.Join(string(config.Directory), dashboard.Folder)

		err = eventuallyCreateDirectory(dashboardFolderInFs)
		if err != nil {
			panic(err)
		}

		// update the file with the new dashboard.
	}
}
