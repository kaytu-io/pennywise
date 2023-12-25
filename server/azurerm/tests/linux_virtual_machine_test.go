package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestLinuxVirtualMachine() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machine data ingested")

	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage data ingested")

	usg, err := ts.getUsage("../../testdata/azure/linux_virtual_machine/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/linux_virtual_machine", *usg)
	costComponent := stat.GetCostComponents()
	expectedCostComponent := []cost.Component{
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.099),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5.28),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.NewFromFloat(0),
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
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_B1s",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
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
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_B1s",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_A2_v2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.091),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_A2_v2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.091),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(76.8),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromFloat(1),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5.28),
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
}
