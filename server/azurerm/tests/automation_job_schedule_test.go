package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestAutomationJobSchedule() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Automation", "eastus")
	fmt.Println("Automation ingested")

	usg, err := ts.getUsage("../../testdata/azure/automation_job_schedule/usage.yaml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/automation_job_schedule", *usg)

	costComponent := stat.GetCostComponents()

	expectedCostComponent := []cost.Component{
		{
			Name:            "Job run time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "minutes",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Job run time",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "minutes",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.002),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Job run time",
			MonthlyQuantity: decimal.NewFromFloat32(5),
			HourlyQuantity:  decimal.Zero,
			Unit:            "minutes",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.002),
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
