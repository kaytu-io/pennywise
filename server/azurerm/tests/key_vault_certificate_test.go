package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestKeyVaultCertificate() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Key Vault", "eastus")
	fmt.Println("Key Vault ingested")

	usg, err := ts.getUsage("../../testdata/azure/key_vault_certificate/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/key_vault_certificate", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Certificate renewals",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(3),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Certificate operations",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.03),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Certificate renewals",
			MonthlyQuantity: decimal.NewFromInt(100),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(3),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Certificate operations",
			MonthlyQuantity: decimal.NewFromFloat(0.01),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.03),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Certificate renewals",
			MonthlyQuantity: decimal.NewFromInt(100),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "requests",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(3),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Certificate operations",
			MonthlyQuantity: decimal.NewFromFloat(0.01),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.03),
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
