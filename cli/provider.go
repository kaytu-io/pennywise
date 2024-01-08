package main

import (
	"github.com/kaytu-io/pennywise/cli/parser/aws"
	"github.com/kaytu-io/pennywise/cli/parser/azurerm"
	"github.com/kaytu-io/pennywise/cli/parser/terraform"
)

// defaultProviders are the currently known and supported terraform providers
var defaultProviders = []terraform.ProviderInitializer{
	aws.TerraformProviderInitializer,
	azurerm.TerraformProviderInitializer,
}
