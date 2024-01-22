package terraform

import (
	"github.com/kaytu-io/pennywise/pkg/cost"
	"github.com/kaytu-io/pennywise/pkg/parser/aws"
	"github.com/kaytu-io/pennywise/pkg/parser/azurerm"
	terraform2 "github.com/kaytu-io/pennywise/pkg/parser/terraform"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/kaytu-io/pennywise/pkg/submission"
	"github.com/kaytu-io/pennywise/pkg/usage"
	"io"
	"os"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

// EstimateTerraformPlanJson is a helper function that reads a Terraform plan json file using the provided io.Reader,
// calculates the costs of the resources and show them.
// It uses the Backend to retrieve the pricing data.
func EstimateTerraformPlanJson(plan io.Reader, u usage.Usage) (*cost.State, error) {
	providerInitializers := []terraform2.ProviderInitializer{
		aws.TerraformProviderInitializer,
		azurerm.TerraformProviderInitializer,
	}

	tfplan := terraform2.NewPlan(providerInitializers...)
	if err := tfplan.Read(plan); err != nil {
		return nil, err
	}
	tfplan.SetUsage(u)

	plannedQueries, err := tfplan.ExtractPlannedQueries()
	if err != nil {
		return nil, err
	}
	var resources []schema.ResourceDef
	for _, rs := range plannedQueries {
		res := rs.ToResource("")
		resources = append(resources, res)
	}
	serverClient := server.NewPennywiseServerClient(ServerClientAddress)
	sub := submission.Submission{Resources: resources}
	stateCost, err := serverClient.GetStateCost(sub)
	if err != nil {
		return nil, err
	}
	return stateCost, nil
}
