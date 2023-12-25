package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestSqlManagedInstance() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("SQL Managed Instance", "westeurope")
	fmt.Println("SQL Managed Instance ingested")

	usg, err := ts.getUsage("../../testdata/azure/sql_managed_instance/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/sql_managed_instance", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Compute (GP_GEN5 16 Cores)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(2.679024),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Additional Storage",
			MonthlyQuantity: decimal.NewFromInt(32),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.13685),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "PITR backup storage (LRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.119),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "LTR backup storage (LRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0298),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "LTR backup storage (LRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0298),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Additional Storage",
			MonthlyQuantity: decimal.NewFromInt(32),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.13685),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "PITR backup storage (ZRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.149),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "LTR backup storage (ZRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0372),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "PITR backup storage (LRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.119),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "SQL license",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(4),
			Unit:            "vCore-hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.099966),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "LTR backup storage (LRS)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.0298),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute (GP_GEN5 4 Cores)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.669756),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Additional Storage",
			MonthlyQuantity: decimal.NewFromInt(32),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.13685),
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
