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
	MySQLPort     = os.Getenv("MYSQL_PORT")
	MySQLDb       = os.Getenv("MYSQL_DB")
	MySQLUser     = os.Getenv("MYSQL_USERNAME")
	MySQLPassword = os.Getenv("MYSQL_PASSWORD")
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
	provider, err := resources2.NewProvider(azurerm.ProviderName)
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

//	func (ts *AzureTestSuite) TestLoadBalancer() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Load Balancer", "")
//		fmt.Println("Load Balancer data ingested")
//
//		lbUsage := usage.Usage{"azurerm_lb": map[string]interface{}{
//			"monthly_data_proceed": 1000,
//		}}
//		cost := ts.getDirCosts("../../testdata/azure/load_balancer", lbUsage)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPublicIp() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Network", "")
//		fmt.Println("Virtual Network data ingested")
//
//		usg := usage.Usage{}
//		cost := ts.getDirCosts("../../testdata/azure/public_ip", usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPublicIpPrefix() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Network", "")
//		fmt.Println("Virtual Network data ingested")
//
//		usg := usage.Usage{}
//		cost := ts.getDirCosts("../../testdata/azure/public_ip_prefix", usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateEndpoint() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Network", "")
//		fmt.Println("Virtual Network data ingested")
//
//		usg := usage.Usage{"azurerm_private_endpoint": map[string]interface{}{
//			"monthly_inbound_data_processed_gb":  100,
//			"monthly_outbound_data_processed_gb": 100,
//		}}
//		cost := ts.getDirCosts("../../testdata/azure/private_endpoint", usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestLinuxVirtualMachine() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		//ts.IngestService("Virtual Machines", "eastus")
//		//fmt.Println("Virtual Machine data ingested")
//		//
//		//ts.IngestService("Storage", "eastus")
//		//fmt.Println("Storage data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/linux_virtual_machine/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/linux_virtual_machine", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestWindowsVirtualMachine() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Machines", "eastus")
//		fmt.Println("Virtual Machine data ingested")
//
//		ts.IngestService("Storage", "eastus")
//		fmt.Println("Storage data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/windows_virtual_machine/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/windows_virtual_machine", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestWindowsVirtualMachineScaleSet() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Machines", "eastus")
//		fmt.Println("Virtual Machine data ingested")
//
//		ts.IngestService("Storage", "eastus")
//		fmt.Println("Storage data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/windows_virtual_machine_scale_set/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/windows_virtual_machine_scale_set", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestVirtualMachine() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Machines", "eastus")
//		fmt.Println("Virtual Machines data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/virtual_machine/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/virtual_machine", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestVirtualMachineScaleSet() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Machines", "eastus")
//		fmt.Println("Virtual Machine data ingested")
//
//		ts.IngestService("Storage", "eastus")
//		fmt.Println("Storage data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/virtual_machine_scale_set/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/virtual_machine_scale_set", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestVirtualNetworkGateway() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Machines", "eastus")
//		fmt.Println("Virtual Machine data ingested")
//
//		ts.IngestService("VPN Gateway", "eastus")
//		fmt.Println("VPN Gateway ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/virtual_network_gateway/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/virtual_network_gateway", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestVirtualNetworkGatewayConnection() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//
//		ts.IngestService("VPN Gateway", "eastus")
//		fmt.Println("VPN Gateway data ingested")
//
//		cost := ts.getDirCosts("../../testdata/azure/virtual_network_gateway_connection", nil)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestVirtualNetworkPeering() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("VPN Gateway", "eastus")
//		fmt.Println("VPN Gateway ingested")
//
//		ts.IngestService("Virtual Network", "eastus")
//		fmt.Println("Virtual Network ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/virtual_network_peering/usage.yml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/virtual_network_peering", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSARecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_a_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_a_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSAAAARecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_aaaa_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_aaaa_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSCAARecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_caa_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_caa_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSCNAMERecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_cname_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_cname_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSMXRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_mx_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_mx_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSNSRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_ns_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_ns_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSPTRRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_ptr_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_ptr_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSSRVRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_srv_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_srv_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSTXTRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/dns_txt_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_txt_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestDNSZoneRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		cost := ts.getDirCosts("../../testdata/azure/dns_zone_record", nil)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSARecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_a_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_a_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSAAAARecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_aaaa_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_aaaa_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSCNAMERecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_cname_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_cname_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSMXRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_mx_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_mx_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSPTRRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_ptr_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_ptr_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSSRVRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_srv_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_srv_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSTXTRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_dns_txt_record/usage.json")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_txt_record", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateDNSZoneRecord() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Azure DNS", "eastus")
//		fmt.Println("Azure DNS data ingested")
//
//		cost := ts.getDirCosts("../../testdata/azure/private_dns_zone_record", nil)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestPrivateEndpoint() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Virtual Network", "eastus")
//		fmt.Println("Virtual Network ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/private_endpoint/usage.yaml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/private_endpoint", *usg)
//		fmt.Println(cost.CostString())
//	}
//
//	func (ts *AzureTestSuite) TestApplicationGateway() {
//		ts.SetupSuite()
//		fmt.Println("Suite Setup")
//		ts.IngestService("Application Gateway", "eastus")
//		fmt.Println("Application Gateway data ingested")
//
//		usg, err := ts.getUsage("../../testdata/azure/application_gateway/usage.yaml")
//		require.NoError(ts.T(), err)
//
//		cost := ts.getDirCosts("../../testdata/azure/application_gateway", *usg)
//		fmt.Println(cost.CostString())
//	}
func (ts *AzureTestSuite) TestCDNEndpoint() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	//ts.IngestService("Content Delivery Network", "")
	//fmt.Println("Content Delivery Network data ingested")

	usg, err := ts.getUsage("../../testdata/azure/cdn_endpoint/usage.yaml")
	require.NoError(ts.T(), err)

	cost := ts.getDirCosts("../../testdata/azure/cdn_endpoint", *usg)
	fmt.Println(cost.GetCostComponents())
}

//
//func (ts *AzureTestSuite) TestContainerRegistry() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Container Registry", "eastus")
//	fmt.Println("Container Registry data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/container_registry/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/container_registry", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbCassandraKeyspace() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_cassandra_keyspace/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_cassandra_keyspace", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbCassandraTable() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_cassandra_table/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/azure/cosmosdb_cassandra_table", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestCosmosdbGremlinDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB data ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_gremlin_database/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_gremlin_graph/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_mongo_collection/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_mongo_database/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_sql_container/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_sql_database/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/cosmosdb_table/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/image/usage.yaml")
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
//	usg, err := ts.getUsage("../../testdata/azure/key_vault_certificate/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/key_vault_certificate", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestKeyVaultKey() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Key Vault", "eastus")
//	fmt.Println("Key Vault ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/key_vault_key/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/key_vault_key", *usg)
//	fmt.Println(cost.CostString())
//}
//func (ts *AzureTestSuite) TestKeyVaultManagedHardwareSecurityModule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Key Vault", "eastus")
//	fmt.Println("Key Vault ingested")
//
//	cost := ts.getDirCosts("../../testdata/key_vault_managed_hardware_security_module", nil)
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
//	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/kubernetes_cluster", *usg)
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
//	usg, err := ts.getUsage("../../testdata/azure/kubernetes_cluster_node_pool/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/kubernetes_cluster_node_pool", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/load_balancer/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/load_balancer", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancerRule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	cost := ts.getDirCosts("../../testdata/lb_rule", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestLoadBalancerOutboundRule() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Load Balancer", "eastus")
//	fmt.Println("Load Balancer ingested")
//
//	cost := ts.getDirCosts("../../testdata/lb_outbound_rule", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestManagedDisk() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/managed_disk/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/managed_disk", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMariadbServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mariadb_server/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/mariadb_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMssqlDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Database", "eastus")
//	fmt.Println("SQL Database ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mssql_database/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/mssql_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMssqlManagedInstance() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Database", "eastus")
//	fmt.Println("SQL Database ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mssql_managed_instance/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/mssql_managed_instance", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMysqlFlexibleServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Database for MySQL", "eastus")
//	fmt.Println("Azure Database for MySQL ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mysql_flexible_server/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/mysql_flexible_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestMysqlServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/mysql_server/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/mysql_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestNatGateway() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("NAT Gateway", "eastus")
//	fmt.Println("NAT Gateway ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/nat_gateway/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/nat_gateway", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPostgresqlFlexibleServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Database for PostgreSQL", "eastus")
//	fmt.Println("Azure Database for PostgreSQL ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/postgresql_flexible_server/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/postgresql_flexible_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPostgresqlServer() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Databases", "eastus")
//	fmt.Println("Databases ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/postgresql_server/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/postgresql_server", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPublicIp() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Network", "eastus")
//	fmt.Println("Virtual Network ingested")
//
//	cost := ts.getDirCosts("../../testdata/public_ip", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestPublicIpPrefix() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Virtual Network", "eastus")
//	fmt.Println("Virtual Network ingested")
//
//	cost := ts.getDirCosts("../../testdata/public_ip_prefix", nil)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSnapshot() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/snapshot/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/snapshot", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSqlDatabase() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Azure Cosmos DB", "eastus")
//	fmt.Println("Azure Cosmos DB ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/sql_database/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/sql_database", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestSqlManagedInstance() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("SQL Managed Instance", "eastus")
//	fmt.Println("SQL Managed Instance ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/sql_managed_instance/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/sql_managed_instance", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageAccount() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_account/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/storage_account", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageQueue() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_queue/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/storage_queue", *usg)
//	fmt.Println(cost.CostString())
//}
//
//func (ts *AzureTestSuite) TestStorageShare() {
//	ts.SetupSuite()
//	fmt.Println("Suite Setup")
//	ts.IngestService("Storage", "eastus")
//	fmt.Println("Storage ingested")
//
//	usg, err := ts.getUsage("../../testdata/azure/storage_share/usage.yaml")
//	require.NoError(ts.T(), err)
//
//	cost := ts.getDirCosts("../../testdata/storage_share", *usg)
//	fmt.Println(cost.CostString())
//}
