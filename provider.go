package main

import (
	"github.com/kaytu-io/pennywise/pkg/parser/aws"
	"github.com/kaytu-io/pennywise/pkg/parser/azurerm"
	"github.com/kaytu-io/pennywise/pkg/parser/terraform"
)

// defaultProviders are the currently known and supported terraform providers
var defaultProviders = []terraform.ProviderInitializer{
	aws.TerraformProviderInitializer,
	azurerm.TerraformProviderInitializer,
}
