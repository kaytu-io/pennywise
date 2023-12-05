package terraform

import (
	"github.com/kaytu-io/pennywise/cli/parser/aws"
	"github.com/kaytu-io/pennywise/cli/parser/azurerm"
	"github.com/kaytu-io/pennywise/cli/parser/terraform"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/client"
	"io"
	"os"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

// EstimateTerraformPlanJson is a helper function that reads a Terraform plan json file using the provided io.Reader,
// calculates the costs of the resources and show them.
// It uses the Backend to retrieve the pricing data.
func EstimateTerraformPlanJson(plan io.Reader, u usage.Usage) error {
	providerInitializers := getDefaultProviders()

	tfplan := terraform.NewPlan(providerInitializers...)
	if err := tfplan.Read(plan); err != nil {
		return err
	}
	tfplan.SetUsage(u)

	plannedQueries, err := tfplan.ExtractPlannedQueries()
	if err != nil {
		return err
	}
	for _, rs := range plannedQueries {
		res := rs.ToResource("")
		serverClient := client.NewPennywiseServerClient(ServerClientAddress)
		_, err := serverClient.GetCost(res)
		if err != nil {
			return err
		}
		//fmt.Println(cost.CostString())
	}
	return nil
}

// defaultProviders are the currently known and supported terraform providers
var defaultProviders = []terraform.ProviderInitializer{
	aws.TerraformProviderInitializer,
	azurerm.TerraformProviderInitializer,
}

// getDefaultProviders will return the default supported providers of terracost
func getDefaultProviders() []terraform.ProviderInitializer {
	return defaultProviders
}
