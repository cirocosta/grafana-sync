package main

import (
	"context"
	"fmt"

	"github.com/cirocosta/grafana-sync/grafana"
	"github.com/jessevdk/go-flags"
)

var config struct {
	AccessToken string        `long:"access-token" required:"true" description:"access token to authenticate against grafana"`
	Address     string        `long:"address" default:"http://localhost:3000" description:"grafana address"`
	Directory   DirectoryFlag `long:"directory" default:"./" description:"directory where dashboards live"`
	Verbose     bool          `short:"v" long:"verbose" description:"displays requests on stderr"`
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

	// perform some validations:
	// 1. does the directory exist?

	for _, dashboard := range dashboards {
		fmt.Printf("%+v\n", dashboard)
	}
}
