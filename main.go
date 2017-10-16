package main

import (
	"github.com/ewilde/terraform-provider-runscope/plugin/providers/runscope"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: runscope.Provider,
	})
}
