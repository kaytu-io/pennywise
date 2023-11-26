package main

import (
	"github.com/kaytu.io/pennywise/server/aws"
	"github.com/kaytu.io/pennywise/server/azurerm"
	"github.com/kaytu.io/pennywise/server/google"
	"github.com/kaytu.io/pennywise/server/terraform"
)

// defaultProviders are the currently known and supported terraform providers
var defaultProviders = []terraform.ProviderInitializer{
	aws.TerraformProviderInitializer,
	google.TerraformProviderInitializer,
	azurerm.TerraformProviderInitializer,
}

// getDefaultProviders will return the default supported providers of terracost
func getDefaultProviders() []terraform.ProviderInitializer {
	return defaultProviders
}
