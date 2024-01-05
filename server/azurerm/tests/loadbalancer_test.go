package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestLoadBalancer() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Load Balancer", "")
	fmt.Println("Load Balancer ingested")

	usg, err := ts.getUsage("../../testdata/azure/load_balancer/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/load_balancer", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Regional Data Proceed",
			MonthlyQuantity: decimal.NewFromInt(100),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Regional Data Proceed",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
	}

	ts.Equal(len(expectedCostComponents), len(costComponents))
	for _, comp := range expectedCostComponents {
		ts.True(componentExists(comp, costComponents), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
