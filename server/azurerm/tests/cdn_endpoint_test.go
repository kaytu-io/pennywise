package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestCDNEndpoint() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Content Delivery Network", "Zone 1")
	fmt.Println("Content Delivery Network data ingested")

	usg, err := ts.getUsage("../../testdata/azure/cdn_endpoint/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/cdn_endpoint", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Outbound data transfer (Premium Verizon, next 350TB)",
			MonthlyQuantity: decimal.NewFromFloat(350000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.102),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Premium Verizon, next 500TB)",
			MonthlyQuantity: decimal.NewFromFloat(500000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.093),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Premium Verizon, first 10TB)",
			MonthlyQuantity: decimal.NewFromFloat(10000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.158),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Premium Verizon, next 40TB)",
			MonthlyQuantity: decimal.NewFromFloat(40000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.14),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Premium Verizon, next 100TB)",
			MonthlyQuantity: decimal.NewFromFloat32(100000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.121),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (next 100TB)",
			MonthlyQuantity: decimal.NewFromFloat(100000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.158),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (next 500TB)",
			MonthlyQuantity: decimal.NewFromFloat(500000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.121),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (over 1000TB)",
			MonthlyQuantity: decimal.NewFromFloat(1000000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.102),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, first 10TB)",
			MonthlyQuantity: decimal.NewFromFloat(10000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.081),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, next 500TB)",
			MonthlyQuantity: decimal.NewFromFloat(500000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.028),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, over 5000TB)",
			MonthlyQuantity: decimal.NewFromFloat(5000000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.023),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (first 50TB)",
			MonthlyQuantity: decimal.NewFromFloat(50000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.177),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (next 350TB)",
			MonthlyQuantity: decimal.NewFromFloat(350000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.14),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, next 40TB)",
			MonthlyQuantity: decimal.NewFromFloat32(40000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.075),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, next 100TB)",
			MonthlyQuantity: decimal.NewFromFloat(100000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.056),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, next 350TB)",
			MonthlyQuantity: decimal.NewFromFloat(350000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.037),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Standard Verizon, next 4000TB)",
			MonthlyQuantity: decimal.NewFromFloat(4000000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.023),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Acceleration outbound data transfer (first 50TB)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.177),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, first 10TB)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.081),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, next 350TB)",
			MonthlyQuantity: decimal.NewFromFloat(350000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.037),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, next 500TB)",
			MonthlyQuantity: decimal.NewFromFloat(200000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.028),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, first 10TB)",
			MonthlyQuantity: decimal.NewFromFloat(10000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.081),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, next 40TB)",
			MonthlyQuantity: decimal.NewFromFloat(40000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.075),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Akamai, next 100TB)",
			MonthlyQuantity: decimal.NewFromFloat(100000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.056),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer (Microsoft, first 10TB)",
			MonthlyQuantity: decimal.NewFromFloat(10000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.081),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Rules engine rules (over 5)",
			MonthlyQuantity: decimal.NewFromFloat(3),
			HourlyQuantity:  decimal.Zero,
			Unit:            "rules",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Rules engine requests",
			MonthlyQuantity: decimal.NewFromFloat(10),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.6),
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
