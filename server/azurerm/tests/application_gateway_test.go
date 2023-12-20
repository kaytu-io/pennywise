package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestApplicationGateway() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Application Gateway", "eastus")
	fmt.Println("Application Gateway data ingested")

	ts.IngestService("Virtual Network", "eastus")
	fmt.Println("Virtual Network data ingested")

	usg, err := ts.getUsage("../../testdata/azure/application_gateway/usage.yml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/application_gateway", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Gateway usage (Basic, Small)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(2),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.025),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (0-10TB)",
			MonthlyQuantity: decimal.NewFromFloat(10240),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (10-40TB)",
			MonthlyQuantity: decimal.NewFromFloat(30720),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (over 40TB)",
			MonthlyQuantity: decimal.NewFromFloat(59040),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF, Medium)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat32(2),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.126),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (0-10TB)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF v2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.36),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "V2 capacity units (WAF)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(2),
			Unit:            "CU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.014),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF, Medium)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(20),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.126),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (10-40TB)",
			MonthlyQuantity: decimal.NewFromFloat(30720),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.007),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (over 40TB)",
			MonthlyQuantity: decimal.NewFromFloat(59040),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.007),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF, Medium)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(2),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.126),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (10-40TB)",
			MonthlyQuantity: decimal.NewFromFloat(30720),
			HourlyQuantity:  decimal.NewFromFloat(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.007),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (over 40TB)",
			MonthlyQuantity: decimal.NewFromFloat32(59040),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.007),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF, Medium)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat32(3),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.126),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Data processing (0-10TB)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF v2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.36),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "V2 capacity units (WAF)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(3),
			Unit:            "CU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.014),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF v2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.36),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "V2 capacity units (WAF)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(20),
			Unit:            "CU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.014),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "IP address (dynamic)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.004),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Gateway usage (WAF v2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.36),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "V2 capacity units (WAF)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(2),
			Unit:            "CU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.014),
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
