package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestKubernetesClusterNodePool() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Machines", "eastus")
	fmt.Println("Virtual Machines ingested")

	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster_node_pool/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/kubernetes_cluster_node_pool", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{
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
				Decimal:  decimal.NewFromFloat(5.888),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
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
				Decimal:  decimal.NewFromFloat(5.888),
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
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.133),
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
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromInt(450),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Basic_A2",
			MonthlyQuantity: decimal.NewFromInt(450),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.079),
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
			Name:            "Compute Standard_DS2_v2",
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
				Decimal:  decimal.NewFromFloat(19.71),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Compute Standard_D2_v3",
			MonthlyQuantity: decimal.NewFromInt(730),
			HourlyQuantity:  decimal.NewFromInt(0),
			Unit:            "Monthly Hours",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.096),
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
	}

	ts.Equal(len(expectedCostComponents), len(costComponents))
	for _, comp := range expectedCostComponents {
		ts.True(componentExists(comp, costComponents), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
}
