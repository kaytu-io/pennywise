package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestVirtualNetworkGateway() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machine data ingested")

	ts.IngestService("VPN Gateway", "Zone 1")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("VPN Gateway", "eastus")
	fmt.Println("VPN Gateway ingested")

	ts.IngestService("Virtual Network", "eastus")
	fmt.Println("Virtual Network data ingested")

	usg, err := ts.getUsage("../../testdata/azure/virtual_network_gateway/usage.yml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/virtual_network_gateway", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "VPN gateway (VpnGw2)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.49),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway P2S tunnels (over 128)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(22),
			Unit:            "tunnel",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
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
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.NewFromFloat(20),
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
			Name:            "VPN gateway (VpnGw5)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(3.65),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway P2S tunnels (over 128)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(5872),
			Unit:            "tunnel",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "VPN gateway P2S tunnels (over 128)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(22),
			Unit:            "tunnel",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "VPN gateway (VpnGw1AZ)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.361),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway (VpnGw3)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "VPN gateway (Basic)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.036),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway P2S tunnels (over 128)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(22),
			Unit:            "tunnel",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.01),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.NewFromFloat32(1),
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
			Name:            "VPN gateway (VpnGw1)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.19),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "VPN gateway data tranfer",
			MonthlyQuantity: decimal.Zero,
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
			Name:            "VPN gateway (VpnGw4)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.1),
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
