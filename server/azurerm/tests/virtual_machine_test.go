package tests

//
//import (
//	"fmt"
//	"github.com/kaytu-io/pennywise/server/cost"
//	"github.com/shopspring/decimal"
//	"github.com/stretchr/testify/require"
//)
//
//func (ts *AzureTestSuite) TestVirtualMachine() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Machines", "eastus")
//	fmt.Println("Virtual Machines data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/virtual_machine/usage.yml")
//	require.NoError(ts.T(), err)
//
//	stat := ts.getDirCosts("../../testdata/azure/virtual_machine", *usg)
//	costComponent := stat.GetCostComponents()
//	expectedCostComponent := []cost.Component{
//		{
//			Name: "Compute Standard_DS1_v2",
//			Unit: "Monthly Hours",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.073),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra disk reservation (if unattached)",
//			Unit: "vCPU",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.006),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(1),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Compute Standard_DS1_v2",
//			Unit: "Monthly Hours",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.126),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra disk reservation (if unattached)",
//			Unit: "vCPU",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.006),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(1),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(2.4),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(5.28),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra LRS Throughput",
//			Rate: cost.Cost{
//				Decimal:  decimal.Zero,
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(8),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra LRS Capacity",
//			Rate: cost.Cost{
//				Decimal:  decimal.Zero,
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(1024),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra LRS IOPs",
//			Rate: cost.Cost{
//				Decimal:  decimal.Zero,
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(2048),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Compute Standard_DS1_v2",
//			Unit: "Monthly Hours",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.126),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra disk reservation (if unattached)",
//			Unit: "vCPU",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.006),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(1),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Compute Standard_DS1_v2",
//			Unit: "Monthly Hours",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.073),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(730),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Ultra disk reservation (if unattached)",
//			Unit: "vCPU",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(0.006),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.Zero,
//			HourlyQuantity:  decimal.NewFromFloat(1),
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//		{
//			Name: "Managed Storage",
//			Rate: cost.Cost{
//				Decimal:  decimal.NewFromFloat(1.536),
//				Currency: "USD",
//			},
//			MonthlyQuantity: decimal.NewFromFloat(1),
//			HourlyQuantity:  decimal.Zero,
//			Usage:           false,
//			Details:         nil,
//			Error:           nil,
//		},
//	}
//
//	ts.Equal(len(expectedCostComponent), len(costComponent))
//	for _, comp := range expectedCostComponent {
//		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
//	}
//}
