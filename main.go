package main

import (
	"github.com/rancher/plugin-server/pkg/filewatcher"
	"github.com/rancher/plugin-server/pkg/server"
	cli "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Start struct {
	Dir   string `usage:"Provide the plugin directory to serve from."`
	Debug bool   `usage:"Set this field to 1 to enable debug logs."`
}

func (a *Start) Run(cmd *cobra.Command, args []string) error {
	//Checks CLI input
	if len(args) == 0 {
		return cmd.Help()
	}

	//Init Logrus
	if a.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	//Registers files in the files directory and starts filewatching service
	logrus.Infof("Creating FileWatcher")
	fw := new(filewatcher.FileWatcher)
	fw.Init(a.Dir)
	//Start filewatcher
	logrus.Infof("Starting FileWatcher")
	go fw.Start()

	//Creates HTTP Fileserver and appropriate handlers
	logrus.Infof("Creating FileServer")
	srv := new(server.FileServer)
	srv.Init(a.Dir, fw)

	//Starts Server
	logrus.Infof("Starting FileServer")
	logrus.Fatal(srv.Srv.ListenAndServe())
	return nil
}

func main() {
	cmd := cli.Command(&Start{}, cobra.Command{
		Long: "Add some long description",
	})

	cli.Main(cmd)
}
