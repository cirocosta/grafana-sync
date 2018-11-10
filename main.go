package main

import (
	_ "github.com/cirocosta/grafana-sync/grafana"
	"github.com/jessevdk/go-flags"
)

var config struct {
	AccessToken string `long:"access-token" required:"true" description:"access token to authenticate against grafana"`
	Address     string `long:"address" default:"http://localhost:3000" description:"grafana address"`
	Directory   string `long:"directory" default:"./" description:"directory where dashboards live"`
}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		return
	}
}
