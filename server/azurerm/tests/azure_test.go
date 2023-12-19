package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kaytu-io/pennywise/cli/parser/hcl"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/azurerm"
	resources2 "github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/cost"
	ingester2 "github.com/kaytu-io/pennywise/server/internal/ingester"
	"github.com/kaytu-io/pennywise/server/internal/mysql"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"testing"
)

var (
	MySQLHost     = os.Getenv("MYSQL_HOST")
	MySQLPort     = "3306"
	MySQLDb       = "test_terracost2"
	MySQLUser     = "test-cost-estimator"
	MySQLPassword = "password"
)

type AzureTestSuite struct {
	suite.Suite

	backend *mysql.Backend
}

func TestAzure(t *testing.T) {
	suite.Run(t, &AzureTestSuite{})
}

func (ts *AzureTestSuite) SetupSuite() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", MySQLUser, MySQLPassword, MySQLHost, MySQLPort, MySQLDb)
	db, err := sql.Open("mysql", dataSource)
	require.NoError(ts.T(), err)
	err = mysql.Migrate(context.Background(), db, "pricing_migrations")

	ts.backend = mysql.NewBackend(db)
}

func (ts *AzureTestSuite) IngestService(service, region string) {
	ingester, err := azurerm.NewIngester(service, region)
	require.NoError(ts.T(), err)

	err = ingester2.IngestPricing(context.Background(), ts.backend, ingester)
	require.NoError(ts.T(), err)

}

func (ts *AzureTestSuite) getUsage(usagePath string) (*usage.Usage, error) {
	var usg usage.Usage
	if usagePath != "" {
		usageFile, err := os.Open(usagePath)
		if err != nil {
			return nil, fmt.Errorf("error while reading usage file %s", err)
		}
		defer usageFile.Close()

		ext := filepath.Ext(usagePath)
		switch ext {
		case ".json":
			err = json.NewDecoder(usageFile).Decode(&usg)
		case ".yaml", ".yml":
			err = yaml.NewDecoder(usageFile).Decode(&usg)
		default:
			return nil, fmt.Errorf("unsupported file format %s for usage file", ext)
		}
		if err != nil {
			return nil, fmt.Errorf("error while parsing usage file %s", err)
		}

	} else {
		usg = usage.Default
	}
	return &usg, nil
}

func (ts *AzureTestSuite) getDirCosts(projectDir string, usg usage.Usage) *cost.State {
	providerName, hclResources, err := hcl.ParseHclResources(projectDir, usg)
	require.NoError(ts.T(), err)

	var qResources []query.Resource
	resources := make(map[string]resource.Resource)
	provider, err := resources2.NewProvider(resources2.ProviderName)
	require.NoError(ts.T(), err)

	for _, rs := range hclResources {
		res := rs.ToResource(providerName, nil)
		resources[res.Address] = res
	}

	for _, res := range resources {
		components := provider.ResourceComponents(resources, res)
		qResource := query.Resource{
			Address:    res.Address,
			Provider:   res.ProviderName,
			Type:       res.Type,
			Components: components,
		}
		qResources = append(qResources, qResource)
	}

	state, err := cost.NewState(context.Background(), ts.backend, qResources)
	require.NoError(ts.T(), err)

	return state
}

func checkComponents(result, expected cost.Component) bool {
	if result.Name == expected.Name && result.MonthlyQuantity.Equal(expected.MonthlyQuantity) &&
		result.HourlyQuantity.Equal(expected.HourlyQuantity) && result.Unit == expected.Unit && result.Rate.Decimal.Equal(expected.Rate.Decimal) &&
		result.Rate.Currency == expected.Rate.Currency && result.Usage == expected.Usage && result.Error == expected.Error {
		return true
	} else {
		return false
	}
}

func componentExists(component cost.Component, comps []cost.Component) bool {
	for _, comp := range comps {
		if checkComponents(comp, component) {
			return true
		}
	}
	return false
}

//

//

//

//

//
//func (ts *AzureTestSuite) TestCosmosdbGremlinDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_gremlin_database/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_gremlin_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbGremlinGraph() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_gremlin_graph/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_gremlin_graph", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbMongoCollection() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_mongo_collection/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_mongo_collection", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbMongoDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_mongo_database/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_mongo_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbSqlContainer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_sql_container/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_sql_container", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbSqlDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_sql_database/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_sql_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbTable() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_table/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_table", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestImage() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/image/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/image", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestKeyVaultCertificate() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Key Vault", "eastus")
//	fmt.Println("Key Vault ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/key_vault_certificate/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/key_vault_certificate", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestKeyVaultKey() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Key Vault", "eastus")
//	fmt.Println("Key Vault ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/key_vault_key/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/key_vault_key", *usg)
//	fmt.Println(cost.CostString())
//}
//func (ts *AzureTestSuite) TestKeyVaultManagedHardwareSecurityModule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Key Vault", "eastus")
//	fmt.Println("Key Vault ingested")
//
//	cost := ts.getDirCosts("../../testdata/azure/key_vault_managed_hardware_security_module", nil)
//	fmt.Println(cost.CostString())
//}
//func (ts *AzureTestSuite) TestKubernetesCluster() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Kubernetes Service", "eastus")
//	fmt.Println("Azure Kubernetes Service ingested")
//
//	ts.IngestService("Virtual Machines", "eastus")
//	fmt.Println("Virtual Machines ingested")
//
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	ts.IngestService("Azure DNS", "eastus")
//	fmt.Println("Azure DNS ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/kubernetes_cluster", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestKubernetesClusterNodePool() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Machines", "eastus")
//	fmt.Println("Virtual Machines ingested")
//
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster_node_pool/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/kubernetes_cluster_node_pool", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/load_balancer/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/load_balancer", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancerRule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	cost := ts.getDirCosts("../../testdata/azure/lb_rule", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancerOutboundRule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	cost := ts.getDirCosts("../../testdata/azure/lb_outbound_rule", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestManagedDisk() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/managed_disk/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/managed_disk", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMariadbServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mariadb_server/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/mariadb_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMssqlDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Database", "eastus")
//	fmt.Println("SQL Database ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mssql_database/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/mssql_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMssqlManagedInstance() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Database", "eastus")
//	fmt.Println("SQL Database ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mssql_managed_instance/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/mssql_managed_instance", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMysqlFlexibleServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Database for MySQL", "eastus")
//	fmt.Println("Azure Database for MySQL ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mysql_flexible_server/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/mysql_flexible_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMysqlServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mysql_server/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/mysql_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestNatGateway() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("NAT Gateway", "eastus")
//	fmt.Println("NAT Gateway ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/nat_gateway/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/nat_gateway", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPostgresqlFlexibleServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Database for PostgreSQL", "eastus")
//	fmt.Println("Azure Database for PostgreSQL ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/postgresql_flexible_server/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/postgresql_flexible_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPostgresqlServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/postgresql_server/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/postgresql_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPublicIp() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Network", "eastus")
//	fmt.Println("Virtual Network ingested")
//
//	cost := ts.getDirCosts("../../testdata/azure/public_ip", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPublicIpPrefix() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Network", "eastus")
//	fmt.Println("Virtual Network ingested")
//
//	cost := ts.getDirCosts("../../testdata/azure/public_ip_prefix", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSnapshot() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/snapshot/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/snapshot", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSqlDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/sql_database/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/sql_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSqlManagedInstance() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Managed Instance", "eastus")
//	fmt.Println("SQL Managed Instance ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/sql_managed_instance/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/sql_managed_instance", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageAccount() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_account/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/storage_account", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageQueue() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_queue/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/storage_queue", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageShare() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_share/usage.yml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/storage_share", *usg)
//	fmt.Println(cost.CostString())
//}
