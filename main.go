package main

import (
	fw "github.com/rancher/plugin-server/pkg/filewatcher"
	"github.com/rancher/plugin-server/pkg/server"
	cli "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Server struct {
	Dir string `usage:"Provide the plugin directory to serve from."`
}

func (a *Server) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	srv := server.New(a.Dir)
	fw.Start(a.Dir)

	logrus.Fatal(srv.ListenAndServe())
	return nil
}

func main() {
	cmd := cli.Command(&Server{}, cobra.Command{
		Long: "Add some long description",
	})

	cli.Main(cmd)
}
