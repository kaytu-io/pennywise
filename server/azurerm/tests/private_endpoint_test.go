package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// TODO : fix storage account problem because effecting in this resource type as well
func (ts *AzureTestSuite) TestPrivateEndpoint() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Virtual Network", "westus")
	fmt.Println("Virtual Network ingested")
	ts.IngestService("Virtual Network", "Global")
	fmt.Println("Virtual Network ingested")

	ts.IngestService("Storage", "westus")
	fmt.Println("Storage ingested")
	ts.IngestService("Storage", "Global")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/private_endpoint/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/private_endpoint", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
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
			Name:            "Data at rest",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "Snapshots",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "Metadata at rest",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.03),
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
				Decimal:  decimal.NewFromFloat(0.13),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "List operations",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10k operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.072),
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
				Decimal:  decimal.NewFromFloat(0.013),
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
				Decimal:  decimal.NewFromFloat(0.006),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Private endpoint",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "hour",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Private endpoint",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "hour",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data processed (first 1PB)",
			MonthlyQuantity: decimal.NewFromFloat(100),
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
			Name:            "Inbound data processed (first 1PB)",
			MonthlyQuantity: decimal.NewFromFloat(1000000),
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
			Name:            "Private endpoint",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "hour",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data processed (first 1PB)",
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
			Name:            "Outbound data processed (first 1PB)",
			MonthlyQuantity: decimal.NewFromFloat(100),
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
			Name:            "Outbound data processed (first 1PB)",
			MonthlyQuantity: decimal.NewFromFloat(1000000),
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
			Name:            "Private endpoint",
			MonthlyQuantity: decimal.NewFromFloat(730),
			HourlyQuantity:  decimal.Zero,
			Unit:            "hour",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data processed (first 1PB)",
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
			Name:            "Outbound data processed (first 1PB)",
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
			Name:            "Outbound data processed (first 1PB)",
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
	}
	ts.Equal(len(expectedCostComponent), len(costComponent))
	for _, comp := range expectedCostComponent {
		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
