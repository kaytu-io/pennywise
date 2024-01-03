package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
)

func (ts *AzureTestSuite) TestAppServiceCertificateOrder() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Azure App Service", "")
	fmt.Println("Container Registry data ingested")

	state := ts.getDirCosts("../../testdata/azure/app_service_certificate_order", nil)
	costComponent := state.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "SSL certificate (Standard)",
			MonthlyQuantity: decimal.NewFromFloat(0.083333333),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "years",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(69.99),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "SSL certificate (Wildcard)",
			MonthlyQuantity: decimal.NewFromFloat(0.083333333),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "years",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(299.99),
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
