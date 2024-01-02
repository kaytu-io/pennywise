package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestFunctionApp() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Functions", "eastus")
	fmt.Println("Functions ingested")

	usg, err := ts.getUsage("../../testdata/azure/function_app/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/function_app", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.NewFromInt(10000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.NewFromFloat(876180.4425),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.NewFromFloat(354012.3),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.NewFromFloat(876180.4425),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.NewFromFloat(354012.3),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "vCPU (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(4),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.173),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Memory (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(14),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0123),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.NewFromInt(37500),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.NewFromInt(300000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "vCPU (EP1)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.173),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Memory (EP1)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(3.5),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0123),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "All other operations",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.004),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data retrieval",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Blob index",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k tags",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.03),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Early deletion",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0152),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Capacity",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0152),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Write operations",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "List and create container operations",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.05),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Read operations",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "vCPU (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(2),
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.173),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Memory (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(7),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0123),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "vCPU (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "vCPU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.173),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Memory (EP2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0123),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Execution time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB-seconds",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Executions",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.000002),
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
