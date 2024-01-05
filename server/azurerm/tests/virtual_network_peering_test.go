package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestVirtualNetworkPeering() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("VPN Gateway", "Zone 2")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("VPN Gateway", "Zone 1")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "Global")
	fmt.Println("Virtual Network ingested")

	ts.IngestService("Virtual Network", "Zone 1")
	fmt.Println("Virtual Network ingested")

	usg, err := ts.getUsage("../../testdata/azure/virtual_network_peering/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/virtual_network_peering", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Outbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer",
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
			Name:            "Inbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
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
			Name:            "Outbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
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
			Name:            "Inbound data transfer",
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
			Name:            "Outbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Outbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.035),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.09),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Inbound data transfer",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.09),
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
