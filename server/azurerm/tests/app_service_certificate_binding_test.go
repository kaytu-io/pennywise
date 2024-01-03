package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
)

func (ts *AzureTestSuite) TestAppServiceCertificateBinding() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Azure App Service", "")
	fmt.Println("Container Registry data ingested")

	state := ts.getDirCosts("../../testdata/azure/app_service_certificate_binding", nil)
	costComponent := state.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "IP SSL certificate",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.Zero,
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(39),
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
