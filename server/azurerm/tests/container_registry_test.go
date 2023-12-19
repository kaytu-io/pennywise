package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestContainerRegistry() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Container Registry", "westeurope")
	fmt.Println("Container Registry data ingested")

	usg, err := ts.getUsage("../../testdata/azure/container_registry/usage.yml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/container_registry", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Registry usage (Premium)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.667),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.NewFromFloat(540000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Geo replication (1 location)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.667),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Registry usage (Premium)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.667),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.NewFromFloat(540000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Registry usage (Basic)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.167),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Registry usage (Basic)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.167),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage (over 10GB)",
			MonthlyQuantity: decimal.NewFromFloat(140),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.NewFromFloat(540000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Registry usage (Premium)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.667),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.NewFromFloat(540000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Geo replication (1 location)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.667),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Build vCPU",
			MonthlyQuantity: decimal.NewFromFloat(540000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "second",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Registry usage (Basic)",
			MonthlyQuantity: decimal.NewFromFloat(30),
			HourlyQuantity:  decimal.Zero,
			Unit:            "Day",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.167),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage (over 10GB)",
			MonthlyQuantity: decimal.NewFromFloat(540),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
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
