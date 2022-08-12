package main

import (
	"github.com/rancher/plugin-server/pkg/filewatcher"
	server "github.com/rancher/plugin-server/pkg/server"
	cli "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Start struct {
	Dir string `usage:"Provide the plugin directory to serve from."`
}

func (a *Start) Run(cmd *cobra.Command, args []string) error {
	//Checks CLI input
	if len(args) == 0 {
		return cmd.Help()
	}

	//Creates HTTP Fileserver and appropriate handlers
	logrus.Infof("Creating FileServer")
	srv := server.New(a.Dir)

	//Registers files in the files directory and starts filewatching service
	logrus.Infof("Creating FileWatcher")
	fw := new(filewatcher.FileWatcher)
	fw.Init(a.Dir)
	//Start filewatcher
	logrus.Infof("Starting FileWatcher")
	go fw.Start()

	//Starts Server
	logrus.Infof("Starting FileServer")
	logrus.Fatal(srv.ListenAndServe())
	return nil
}

func main() {
	cmd := cli.Command(&Start{}, cobra.Command{
		Long: "Add some long description",
	})

	cli.Main(cmd)
}
