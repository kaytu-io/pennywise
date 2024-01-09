package main

import (
	"github.com/kaytu-io/pennywise/parser/aws"
	"github.com/kaytu-io/pennywise/parser/azurerm"
	"github.com/kaytu-io/pennywise/parser/terraform"
)

// defaultProviders are the currently known and supported terraform providers
var defaultProviders = []terraform.ProviderInitializer{
	aws.TerraformProviderInitializer,
	azurerm.TerraformProviderInitializer,
}
