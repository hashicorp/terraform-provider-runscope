package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ewilde/terraform-provider-runscope/plugin/providers/runscope"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: runscope.Provider,
	})
}
