package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestVirtualMachineScaleSet() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Machines", "westeurope")
	fmt.Println("Virtual Machine data ingested")

	ts.IngestService("Storage", "westeurope")
	fmt.Println("Storage data ingested")

	usg, err := ts.getUsage("../../testdata/azure/virtual_machine_scale_set/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/virtual_machine_scale_set", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.114),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.114),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.114),
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
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.4),
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
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.114),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_F2",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.114),
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
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.4),
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
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
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
