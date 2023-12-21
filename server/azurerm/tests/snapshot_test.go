package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestSnapshot() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/snapshot/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/snapshot", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(70),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "1 GB/Month",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.05),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(20),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "1 GB/Month",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.05),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(70),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "1 GB/Month",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.05),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(20),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "1 GB/Month",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.05),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5.888),
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
