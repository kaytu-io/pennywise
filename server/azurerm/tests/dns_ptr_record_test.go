package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestDNSPTRRecord() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	//ts.IngestService("Azure DNS", "West Europ")
	//fmt.Println("Azure DNS data ingested")

	usg, err := ts.getUsage("../../testdata/azure/dns_ptr_record/usage.json")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/dns_ptr_record", *usg)
	costComponent := stat.GetCostComponents()
	for k, v := range costComponent {
		fmt.Printf("cost component : %v \n", k)
		fmt.Printf("name : %v \n ", v.Name)
		fmt.Printf("unit : %v \n ", v.Unit)
		fmt.Printf("rate : %v \n ", v.Rate)
		fmt.Printf("Details : %v \n ", v.Details)
		fmt.Printf("Usage : %v \n ", v.Usage)
		fmt.Printf("MonthlyQuantity : %v \n ", v.MonthlyQuantity)
		fmt.Printf("HourlyQuantity : %v \n ", v.HourlyQuantity)
		fmt.Printf("Error : %v \n ", v.Error)
		fmt.Printf("\n")
	}
	expectedCostComponent := []cost.Component{
		{
			Name:            "DNS queries (first 1B)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
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
			MonthlyQuantity: decimal.NewFromFloat(1000),
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
			MonthlyQuantity: decimal.NewFromFloat(500),
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
	}
	ts.Equal(len(expectedCostComponent), len(costComponent))
	for _, comp := range expectedCostComponent {
		ts.True(componentExists(comp, costComponent), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
