package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestWindowsVirtualMachineScaleSet() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machine data ingested")

	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage data ingested")

	usg, err := ts.getUsage("../../testdata/azure/windows_virtual_machine_scale_set/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/windows_virtual_machine_scale_set", *usg)
	costComponent := stat.GetCostComponents()
	expectedCostComponent := []cost.Component{
		{
			Name: "Compute Basic_A2",
			Unit: "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Usage:           false,
			Details:         nil,
			Error:           nil,
		},
		{
			Name:            "Compute Basic_A2",
			Unit:            "Monthly Hours",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Compute Basic_A2",
			Unit:            "Monthly Hours",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Compute Basic_A2",
			Unit:            "Monthly Hours",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Compute Basic_A2",
			Unit:            "Monthly Hours",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
		{
			Name:            "Compute Basic_A2",
			Unit:            "Monthly Hours",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
				Currency: "USD",
			},
			Usage:   false,
			Details: nil,
			Error:   nil,
		},
	}
	ts.Equal(len(expectedCostComponent), len(costComponent))
	for _, comp := range expectedCostComponent {
		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
