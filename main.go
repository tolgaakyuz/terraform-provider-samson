package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-samson/samson"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: samson.Provider})
}
