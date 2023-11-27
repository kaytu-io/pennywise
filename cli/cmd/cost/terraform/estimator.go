package terraform

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/parser/aws"
	"github.com/kaytu-io/pennywise/cli/parser/azurerm"
	"github.com/kaytu-io/pennywise/cli/parser/terraform"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/client"
	"io"
	"os"
	"sort"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

// EstimateTerraformPlan is a helper function that reads a Terraform plan using the provided io.Reader,
// generates the prior and planned cost.State, and then creates a cost.Plan from them that is returned.
// It uses the Backend to retrieve the pricing data.
func EstimateTerraformPlan(plan io.Reader, u usage.Usage, providerInitializers ...terraform.ProviderInitializer) error {
	if len(providerInitializers) == 0 {
		providerInitializers = getDefaultProviders()
	}

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
		cost, err := serverClient.GetCost(res)
		if err != nil {
			return err
		}
		fmt.Println(rs.Address, ":", cost)
	}

	modules := make([]string, 0, 0)
	for k := range tfplan.Configuration.RootModule.ModuleCalls {
		modules = append(modules, k)
	}
	sort.Strings(modules)
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
