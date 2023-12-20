package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestDNSNSRecord() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Azure DNS", "Zone 1")
	fmt.Println("Azure DNS data ingested")

	usg, err := ts.getUsage("../../testdata/azure/dns_ns_record/usage.json")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/dns_ns_record", *usg)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
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
			MonthlyQuantity: decimal.NewFromFloat(1),
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
	}
	ts.Equal(len(expectedCostComponent), len(costComponent))
	for _, comp := range expectedCostComponent {
		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
