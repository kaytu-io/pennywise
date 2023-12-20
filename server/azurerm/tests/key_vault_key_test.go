package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestKeyVaultKey() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Key Vault", "eastus")
	fmt.Println("Key Vault ingested")

	usg, err := ts.getUsage("../../testdata/azure/key_vault_key/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/key_vault_key", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(30),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(300),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(3),
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
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(2),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(20),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(200),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys (next 1250)",
			MonthlyQuantity: decimal.NewFromInt(1250),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.5),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys (next 2500)",
			MonthlyQuantity: decimal.NewFromInt(2500),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.9),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys (over 4000)",
			MonthlyQuantity: decimal.NewFromInt(1000),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.4),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys",
			MonthlyQuantity: decimal.NewFromInt(400),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(4),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(40),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys (first 250)",
			MonthlyQuantity: decimal.NewFromInt(250),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(100),
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
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(1),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(10),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(5),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(50),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(500),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(7),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(70),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(700),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
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
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(6),
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
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(60),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Software-protected keys",
			MonthlyQuantity: decimal.NewFromInt(600),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Storage key rotations",
			MonthlyQuantity: decimal.NewFromInt(40),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "renewals",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys",
			MonthlyQuantity: decimal.NewFromInt(5000),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "months",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "HSM-protected keys",
			MonthlyQuantity: decimal.NewFromInt(400),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "10K transactions",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Secrets operations",
			MonthlyQuantity: decimal.NewFromInt(4),
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
