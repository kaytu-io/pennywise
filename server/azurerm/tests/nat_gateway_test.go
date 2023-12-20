package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestNatGateway() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("NAT Gateway", "westeurope")
	fmt.Println("NAT Gateway ingested")
	ts.IngestService("Virtual Network", "westeurope")
	fmt.Println("Virtual Network ingested")

	usg, err := ts.getUsage("../../testdata/azure/nat_gateway/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/nat_gateway", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "IP prefix",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.006),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "NAT gateway",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.045),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processed",
			MonthlyQuantity: decimal.NewFromInt(10),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.045),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "IP address (static)",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "NAT gateway",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.045),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processed",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.045),
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
