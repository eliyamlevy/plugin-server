package main

import (
	"github.com/rancher/rancher-plugins-server/pkg/app"
	cli "github.com/rancher/wrangler-cli"
)

func main() {
	cli.Main(app.New())
}
