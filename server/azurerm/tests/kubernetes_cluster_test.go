package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestKubernetesCluster() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Azure Kubernetes Service", "eastus")
	fmt.Println("Azure Kubernetes Service ingested")

	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machines ingested")

	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage ingested")

	ts.IngestService("Load Balancer", "eastus")
	fmt.Println("Load Balancer ingested")

	ts.IngestService("Azure DNS", "eastus")
	fmt.Println("Azure DNS ingested")

	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/kubernetes_cluster", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5.888),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Uptime SLA",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(1.536),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Uptime SLA",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Uptime SLA",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_DS2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Hosted zone",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
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
			Name:            "Uptime SLA",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(450),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.146),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Regional Data Proceed",
			MonthlyQuantity: decimal.NewFromInt(100),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Uptime SLA",
			MonthlyQuantity: decimal.NewFromInt(0),
			HourlyQuantity:  decimal.NewFromInt(1),
			Unit:            "hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.1),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.252),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.252),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.252),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.252),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.252),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},

		{
			Name:            "Managed Storage",
			MonthlyQuantity: decimal.NewFromInt(1),
			HourlyQuantity:  decimal.NewFromInt(0),
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(5.888),
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
