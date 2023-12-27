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
	ts.IngestService("Virtual Machines", "westeurope")
	fmt.Println("Virtual Machines ingested")
	ts.IngestService("Storage", "westeurope")
	fmt.Println("Storage ingested")
	ts.IngestService("SQL Database", "westeurope")
	fmt.Println("SQL Database ingested")

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
				Decimal:  decimal.NewFromFloat(2.679),
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
				Decimal:  decimal.NewFromFloat(0.137),
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
				Decimal:  decimal.NewFromFloat(0.03),
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
				Decimal:  decimal.NewFromFloat(0.03),
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
				Decimal:  decimal.NewFromFloat(0.137),
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
				Decimal:  decimal.NewFromFloat(0.037),
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
				Decimal:  decimal.NewFromFloat(0.1),
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
				Decimal:  decimal.NewFromFloat(0.67),
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
				Decimal:  decimal.NewFromFloat(0.137),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute (BC_GEN5 40 Cores)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(13.395),
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
	fmt.Println(costComponents)
}
