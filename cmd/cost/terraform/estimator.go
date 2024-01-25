package terraform

import (
	"github.com/kaytu-io/pennywise/pkg/parser/aws"
	"github.com/kaytu-io/pennywise/pkg/parser/azurerm"
	terraform2 "github.com/kaytu-io/pennywise/pkg/parser/terraform"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/usage"
	"io"
	"os"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

// ParseTerraformPlanJson is a helper function that reads a Terraform plan json file using the provided io.Reader,
// calculates the costs of the resources and show them.
// It uses the Backend to retrieve the pricing data.
func ParseTerraformPlanJson(plan io.Reader, u usage.Usage) ([]schema.ResourceDef, error) {
	providerInitializers := []terraform2.ProviderInitializer{
		aws.TerraformProviderInitializer,
		azurerm.TerraformProviderInitializer,
	}

	tfplan := terraform2.NewPlan(providerInitializers...)
	if err := tfplan.Read(plan); err != nil {
		return nil, err
	}
	tfplan.SetUsage(u)
	var defaultRegion string
	for _, config := range tfplan.Configuration.ProviderConfig {
		for key, value := range config.Expressions {
			if _, ok := value.ConstantValue.(string); ok && key == "region" {
				defaultRegion = value.ConstantValue.(string)
			}
		}
	}

	plannedQueries, err := tfplan.ExtractPlannedQueries()
	if err != nil {
		return nil, err
	}
	var resources []schema.ResourceDef
	for _, rs := range plannedQueries {
		res := rs.ToResource(defaultRegion)
		resources = append(resources, res)
	}
	return resources, nil
}
