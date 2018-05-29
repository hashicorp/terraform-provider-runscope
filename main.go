package main

import (
	"github.com/terraform-providers/terraform-provider-runscope/runscope"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: runscope.Provider,
	})
}
