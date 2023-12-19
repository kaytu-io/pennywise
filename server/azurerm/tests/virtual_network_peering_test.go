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
	ts.IngestService("VPN Gateway", "West Europe")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "West Europe")
	fmt.Println("Virtual Network ingested")

	ts.IngestService("VPN Gateway", "North Europe")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "North Europe")
	fmt.Println("Virtual Network ingested")

	ts.IngestService("VPN Gateway", "Japan West")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "Japan West")
	fmt.Println("Virtual Network ingested")

	ts.IngestService("VPN Gateway", "Global")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "Global")
	fmt.Println("Virtual Network ingested")

	usg, err := ts.getUsage("../../testdata/azure/virtual_network_peering/usage.yml")
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
