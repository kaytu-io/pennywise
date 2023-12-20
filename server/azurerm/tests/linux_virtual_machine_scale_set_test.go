package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestLinuxVirtualMachineScaleSet() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machine data ingested")

	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage data ingested")

	usg, err := ts.getUsage("../../testdata/azure/linux_virtual_machine_scale_set/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/linux_virtual_machine_scale_set", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
	}

	ts.Equal(len(expectedCostComponents), len(costComponents))
	for _, comp := range expectedCostComponents {
		ts.True(componentExists(comp, costComponents))
	}
}
