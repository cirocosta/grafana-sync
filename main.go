package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
)

var config struct {
	Address     string        `short:"a" long:"address" default:"http://localhost:3000" description:"grafana address"`
	Directory   DirectoryFlag `short:"d" long:"directory" default:"./" description:"directory where dashboards live"`
	Verbose     bool          `short:"v" long:"verbose" description:"displays requests on stderr"`
	SyncFolders []string      `long:"sync-folders" short:"f" description:"specific folders to sync"`
	Auth        struct {
		Username    string `short:"u" long:"username" description:"basic auth username"`
		Password    string `short:"p" long:"password" description:"basic auth password"`
		AccessToken string `long:"access-token" description:"access token to authenticate against grafana"`
	} `group:"Authentication"`

	Push pushCommand `command:"push" description:"pushes the dashboards to a grafana instance"`
	Pull pullCommand `command:"pull" description:"pulls the dashboards from a grafana instance"`
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
}

func shouldSyncFolder(syncFolders []string, folder string) bool {
	for _, f := range syncFolders {
		if folder == f {
			return true
		}
	}
	return false
}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(1)
	}
}
