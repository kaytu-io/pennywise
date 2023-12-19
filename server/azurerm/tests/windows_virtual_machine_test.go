package tests

//
//import (
//	"fmt"
//	"github.com/kaytu-io/pennywise/server/cost"
//	"github.com/shopspring/decimal"
//	"github.com/stretchr/testify/require"
//)
//
//func (ts *AzureTestSuite) TestWindowsVirtualMachine() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Machines", "eastus")
//	fmt.Println("Virtual Machine data ingested")
//
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/windows_virtual_machine/usage.yml")
//	require.NoError(ts.T(), err)
//
//	stat := ts.getDirCosts("../../testdata/azure/windows_virtual_machine", *usg)
//	costComponent := stat.GetCostComponents()
//	expectedCostComponent := []cost.Component{
//		{
//			Name:            "Compute Standard_A2_v2",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.136),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(2.4),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Compute Standard_A2_v2",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.136),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(76.8),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(76.8),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Compute Standard_D2_v4",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.096),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(76.8),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Compute Standard_F2",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.192),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Compute Standard_E16-8as_v4",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.744),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Compute Basic_A2",
//			Unit:            "Monthly Hours",
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.133),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//		{
//			Name:            "Managed Storage",
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			Details: nil,
//			Usage:   false,
//			Error:   nil,
//		},
//	}
//
//	ts.Equal(len(expectedCostComponent), len(costComponent))
//	for _, comp := range expectedCostComponent {
//		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
//	}
//}
