package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestImage() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "westeurope")
	fmt.Println("Storage ingested")
	ts.IngestService("Virtual Machines", "westeurope")
	fmt.Println("Virtual Machine data ingested")

	usg, err := ts.getUsage("../../testdata/azure/image/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/image", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D1_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.068),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Ultra disk reservation (if unattached)",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(30),
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
			Name:            "Ultra disk reservation (if unattached)",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
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
				Decimal:  decimal.NewFromFloat(5.807),
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
				Decimal:  decimal.NewFromFloat(5.807),
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
				Decimal:  decimal.NewFromFloat(5.807),
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
				Decimal:  decimal.NewFromFloat(5.807),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS1_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.068),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS1_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.068),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Ultra disk reservation (if unattached)",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
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
		{
			Name:            "Compute Standard_DS1_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.068),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Ultra disk reservation (if unattached)",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
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
				Decimal:  decimal.NewFromFloat(5.807),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage",
			MonthlyQuantity: decimal.NewFromInt(60),
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
	}

	ts.Equal(len(expectedCostComponents), len(costComponents))
	for _, comp := range expectedCostComponents {
		ts.True(componentExists(comp, costComponents), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
