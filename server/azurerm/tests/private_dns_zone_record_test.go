package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
)

func (ts *AzureTestSuite) TestPrivateDNSZoneRecord() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Azure DNS", "westus")
	fmt.Println("Azure DNS data ingested")

	stat := ts.getDirCosts("../../testdata/azure/private_dns_zone_record", nil)
	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
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
