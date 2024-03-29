package main

import (
	"github.com/csechrist/terraform-provider-konnect/konnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: konnect.Provider,
	})
}
