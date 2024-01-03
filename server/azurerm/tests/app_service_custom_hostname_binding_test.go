package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
)

func (ts *AzureTestSuite) TestAppServiceCustomHostnameBinding() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Azure App Service", "eastus2")
	fmt.Println("Container Registry data ingested")

	state := ts.getDirCosts("../../testdata/azure/app_service_custom_hostname_binding", nil)
	costComponent := state.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Stamp fee",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.358),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
	}
	ts.Equal(len(expectedCostComponent), len(costComponent))
	for _, comp := range expectedCostComponent {
		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
	fmt.Println(costComponent)
	fmt.Println(state.CostString())
}
