package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestPrivateDNSARecord() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Azure DNS", "Zone 1")
	fmt.Println("Azure DNS data ingested")

	usg, err := ts.getUsage("../../testdata/azure/private_dns_a_record/usage.json")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/private_dns_a_record", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "DNS queries (first 1B)",
			MonthlyQuantity: decimal.NewFromInt(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M queries",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "DNS queries (first 1B)",
			MonthlyQuantity: decimal.NewFromInt(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M queries",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "DNS queries (over 1B)",
			MonthlyQuantity: decimal.NewFromInt(500),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M queries",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.2),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "DNS queries (first 1B)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M queries",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Hosted zone",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.5),
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
