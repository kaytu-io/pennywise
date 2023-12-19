package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestCosmosdbCassandraTable() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	//ts.IngestService("Azure Cosmos DB", "eastus")
	//fmt.Println("Azure Cosmos DB data ingested")
	//
	//ts.IngestService("Azure Cosmos DB", "westus")
	//fmt.Println("Azure Cosmos DB data ingested")
	//
	//ts.IngestService("Azure Cosmos DB", "centralus")
	//fmt.Println("Azure Cosmos DB data ingested")

	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_cassandra_table/usage.yml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/cosmosdb_cassandra_table", *usg)
	costComponent := stat.GetCostComponents()
	for k, v := range costComponent {
		fmt.Printf("cost component : %v \n", k)
		fmt.Printf("name : %v \n ", v.Name)
		fmt.Printf("unit : %v \n ", v.Unit)
		fmt.Printf("rate : %v \n ", v.Rate)
		fmt.Printf("Details : %v \n ", v.Details)
		fmt.Printf("Usage : %v \n ", v.Usage)
		fmt.Printf("MonthlyQuantity : %v \n ", v.MonthlyQuantity)
		fmt.Printf("HourlyQuantity : %v \n ", v.HourlyQuantity)
		fmt.Printf("Error : %v \n ", v.Error)
		fmt.Printf("\n")
	}
	expectedCostComponent := []cost.Component{
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(5),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(6.25),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical storage (West US)",
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
			Name:            "Analytical write operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.055),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical read operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Continuous backup (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.2),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Continuous backup (Central US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.246),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.NewFromFloat(3000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (autoscale, West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(45),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (autoscale, Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(45),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(10),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(10),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (serverless)",
			MonthlyQuantity: decimal.NewFromFloat(10),
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M RU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.279),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical storage (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
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
			Name:            "Analytical write operations (West US)",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.055),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical read operations (West US)",
			MonthlyQuantity: decimal.NewFromFloat(100),
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Periodic backup (West US)",
			MonthlyQuantity: decimal.NewFromFloat(2000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.NewFromFloat(3000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(10),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(10),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.016),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Continuous backup (West US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.2),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Continuous backup (Central US)",
			MonthlyQuantity: decimal.NewFromFloat(1000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.246),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.NewFromFloat(3000),
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(5),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (provisioned, Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.NewFromFloat(6.25),
			Unit:            "RU/s x 100",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.008),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical storage (West US)",
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
			Name:            "Analytical write operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.055),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical read operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Provisioned throughput (serverless)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "1M RU",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.279),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (Central US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Transactional storage (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.25),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical storage (West US)",
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
			Name:            "Analytical write operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.055),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Analytical read operations (West US)",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "10K operations",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.005),
				Currency: "USD",
			},
			Details: []string{},
			Usage:   false,

			Error: nil,
		},
		{
			Name:            "Restored data",
			MonthlyQuantity: decimal.Zero,
			HourlyQuantity:  decimal.Zero,
			Unit:            "GB",
			Rate: cost.Cost{
				Decimal:  decimal.NewFromFloat(0.15),
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
