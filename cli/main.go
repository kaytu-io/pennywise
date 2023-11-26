package main

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/parser/terraform"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/resource"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"sort"
)

func main() {
	file, err := os.Open("../../terracost_test/image/tfplan.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	usage := usage.Default
	err = EstimateTerraformPlan(context.Background(), file, usage)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

// EstimateTerraformPlan is a helper function that reads a Terraform plan using the provided io.Reader,
// generates the prior and planned cost.State, and then creates a cost.Plan from them that is returned.
// It uses the Backend to retrieve the pricing data.
func EstimateTerraformPlan(ctx context.Context, plan io.Reader, u usage.Usage, providerInitializers ...terraform.ProviderInitializer) error {
	if len(providerInitializers) == 0 {
		providerInitializers = getDefaultProviders()
	}

	tfplan := terraform.NewPlan(providerInitializers...)
	if err := tfplan.Read(plan); err != nil {
		return err
	}
	tfplan.SetUsage(u)

	priorQueries, err := tfplan.ExtractPriorQueries()
	if err != nil {
		return err
	}
	fmt.Println("PRIOR")
	fmt.Println(priorQueries)
	// If it's the first time we run the plan, then we might not have
	// prior queries so we ignore it and move forward
	//prior, err := cost.NewState(ctx, be, priorQueries)
	//if err != nil && err != terraform.ErrNoQueries {
	//	return nil, err
	//}

	plannedQueries, err := tfplan.ExtractPlannedQueries()
	if err != nil {
		return err
	}
	fmt.Println("PLANNED")
	fmt.Println(plannedQueries)
	resource.Resource{}
	//planned, err := cost.NewState(ctx, be, plannedQueries)
	//if err != nil {
	//	return nil, err
	//}

	modules := make([]string, 0, 0)
	for k := range tfplan.Configuration.RootModule.ModuleCalls {
		modules = append(modules, k)
	}
	sort.Strings(modules)

	return nil
}
