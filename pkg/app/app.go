package app

import (
	"log"

	fileserver "github.com/rancher/rancher-plugins-server/pkg/fileserver"
	cli "github.com/rancher/wrangler-cli"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	root := cli.Command(&App{}, cobra.Command{
		Long: "Add some long description",
	})
	return root
}

type App struct {
	Dir string `usage:"Provide the plugin directory to serve from."`
}

func (a *App) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	srv := fileserver.New(a.Dir)

	log.Fatal(srv.ListenAndServe())
	return nil
}
