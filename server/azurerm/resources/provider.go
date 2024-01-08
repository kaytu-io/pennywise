package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/resource"
	"go.uber.org/zap"
)

const ProviderName = "azurerm"

var (
	locationDisplayToName = map[string]string{
		"West US":              "westus",
		"West US 2":            "westus2",
		"East US":              "eastus",
		"Central US":           "centralus",
		"Central US EUAP":      "centraluseuap",
		"South Central US":     "southcentralus",
		"North Central US":     "northcentralus",
		"West Central US":      "westcentralus",
		"East US 2":            "eastus2",
		"East US 2 EUAP":       "eastus2euap",
		"Brazil South":         "brazilsouth",
		"Brazil US":            "brazilus",
		"North Europe":         "northeurope",
		"West Europe":          "westeurope",
		"East Asia":            "eastasia",
		"Southeast Asia":       "southeastasia",
		"Japan West":           "japanwest",
		"Japan East":           "japaneast",
		"Korea Central":        "koreacentral",
		"Korea South":          "koreasouth",
		"South India":          "southindia",
		"West India":           "westindia",
		"Central India":        "centralindia",
		"Australia East":       "australiaeast",
		"Australia Southeast":  "australiasoutheast",
		"Canada Central":       "canadacentral",
		"Canada East":          "canadaeast",
		"UK South":             "uksouth",
		"UK West":              "ukwest",
		"France Central":       "francecentral",
		"France South":         "francesouth",
		"Australia Central":    "australiacentral",
		"Australia Central 2":  "australiacentral2",
		"UAE Central":          "uaecentral",
		"UAE North":            "uaenorth",
		"South Africa North":   "southafricanorth",
		"South Africa West":    "southafricawest",
		"Switzerland North":    "switzerlandnorth",
		"Switzerland West":     "switzerlandwest",
		"Germany North":        "germanynorth",
		"Germany West Central": "germanywestcentral",
		"Norway East":          "norwayeast",
		"Norway West":          "norwaywest",
		"Brazil Southeast":     "brazilsoutheast",
		"West US 3":            "westus3",
		"East US SLV":          "eastusslv",
		"Sweden Central":       "swedencentral",
		"Sweden South":         "swedensouth",
	}
)

func locationNameMapping(location string) string {
	newMapping := make(map[string]string)

	for key, value := range locationDisplayToName {
		newMapping[value] = key
	}

	return newMapping[location]
}

// Provider is an implementation of the resources.Provider, used to extract component queries from
// terraform resources.
type Provider struct {
	key    string
	logger *zap.Logger
}

// NewProvider initializes a new Azure provider with key and region
func NewProvider(key string, logger *zap.Logger) (*Provider, error) {
	return &Provider{
		key:    key,
		logger: logger,
	}, nil
}

// Name returns the Provider's common name.
func (p *Provider) Name() string { return p.key }

// ResourceComponents returns Component queries for a given terraform.Resource.
func (p *Provider) ResourceComponents(rss map[string]resource.ResourceDef, tfRes resource.ResourceDef) []resource.Component {
	switch tfRes.Type {
	case "azurerm_linux_virtual_machine":
		vals, err := decodeLinuxVirtualMachineValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newLinuxVirtualMachine(vals).Components()
	case "azurerm_virtual_machine":
		vals, err := decodeVirtualMachineValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newVirtualMachine(vals).Components()
	case "azurerm_windows_virtual_machine":
		vals, err := decodeWindowsVirtualMachineValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newWindowsVirtualMachine(vals).Components()
	case "azurerm_managed_disk":
		vals, err := decodeManagedStorageValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newManagedStorage(vals).Components()
	case "azurerm_image":
		vals, err := decodeImageValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newImage(vals).Components()
	case "azurerm_snapshot":
		vals, err := decodeSnapshotValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newSnapshot(vals).Components()
	case "azurerm_linux_virtual_machine_scale_set":
		vals, err := decodeLinuxVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newLinuxVirtualMachineScaleSet(vals).Components()
	case "azurerm_windows_virtual_machine_scale_set":
		vals, err := decodeWindowsVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newWindowsVirtualMachineScaleSet(vals).Components()
	case "azurerm_virtual_machine_scale_set":
		vals, err := decodeVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newVirtualMachineScaleSet(vals).Components()
	case "azurerm_lb":
		vals, err := decodeLoadBalancerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newLoadBalancer(vals).Components()
	case "azurerm_lb_rule":
		vals, err := decodeLoadBalancerRuleValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newLoadBalancerRule(vals).Components()
	case "azurerm_lb_outbound_rule":
		vals, err := decodeLoadBalancerRuleValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newLoadBalancerRule(vals).Components()
	case "azurerm_application_gateway":
		vals, err := decodeApplicationGatewayValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newApplicationGateway(vals).Components()
	case "azurerm_nat_gateway":
		vals, err := decodeNatGatewayValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newNatGateway(vals).Components()
	case "azurerm_public_ip":
		vals, err := decodePublicIPValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPublicIP(vals).Components()
	case "azurerm_public_ip_prefix":
		vals, err := decodePublicIPPrefixValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPublicIPPrefix(vals).Components()
	case "azurerm_container_registry":
		vals, err := decodeContainerRegistry(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newContainerRegistry(vals).component()
	case "azurerm_private_endpoint":
		vals, err := decodePrivateEndpointValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateEndpoint(vals).Components()
	case "azurerm_storage_queue":
		vals, err := decodeStorageQueueValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newStorageQueue(vals).Components()
	case "azurerm_storage_share":
		vals, err := decodeStorageShareValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newStorageShare(vals).Components()
	case "azurerm_storage_account":
		vals, err := decodeStorageAccountValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newStorageAccount(vals).Components()
	case "azurerm_virtual_network_gateway":
		vals, err := decodeVirtualNetworkGateway(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newVirtualNetworkGateway(vals).Components()
	case "azurerm_virtual_network_gateway_connection":
		vals, err := decodeVirtualNetworkGatewayConnection(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newVirtualNetworkGatewayConnection(vals).Component()
	case "azurerm_key_vault_key":
		vals, err := decodeKeyVaultKeyValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newKeyVaultKey(vals).Components()
	case "azurerm_key_vault_certificate":
		vals, err := decodeKeyVaultCertificateValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newKeyVaultCertificate(vals).Components()
	case "azurerm_key_vault_managed_hardware_security_module":
		vals, err := decodeKeyVaultManagedHardwareSecurityModuleValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newKeyVaultManagedHardwareSecurityModule(vals).Components()
	case "azurerm_virtual_network_peering":
		vals, err := decodeVirtualNetworkPeeringValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newVirtualNetworkPeering(vals).Components()
	case "azurerm_cdn_endpoint":
		vals, err := decodeCDNEndpoint(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCDNEndpoint(vals).Component()
	case "azurerm_dns_a_record":
		vals, err := decodeDNSARecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSARecord(vals).component()
	case "azurerm_dns_aaaa_record":
		vals, err := decoderDNSAAAARecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSAAAARecord(vals).component()
	case "azurerm_dns_caa_record":
		vals, err := decoderDNSCAARecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSCAARecord(vals).component()
	case "azurerm_dns_cname_record":
		vals, err := decoderDNSCNAMERecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSCNAMERecord(vals).component()
	case "azurerm_dns_mx_record":
		vals, err := decoderDNSMXRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSMXRecord(vals).component()
	case "azurerm_dns_ns_record":
		vals, err := decoderDNSNSRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSNSRecord(vals).component()
	case "azurerm_dns_ptr_record":
		vals, err := decoderDNSPTRRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSPTRRecord(vals).component()
	case "azurerm_dns_srv_record":
		vals, err := decoderDNSSRVRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSSRVRecord(vals).component()
	case "azurerm_dns_txt_record":
		vals, err := decoderDNSTXTRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newDNSTXTRecord(vals).component()
	case "azurerm_dns_zone":
		vals, err := decoderRMDNSZone(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newRMDNSZone(vals).component()
	case "azurerm_private_dns_a_record":
		vals, err := decoderPrivateDnsARecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDnsARecord(vals).component()
	case "azurerm_private_dns_aaaa_record":
		vals, err := decoderPrivateDnsAAAARecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newprivateDNSAAAARecord(vals).component()
	case "azurerm_private_dns_cname_record":
		vals, err := decoderPrivateDnsCNAMERecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newprivateDNSCNAMERecord(vals).component()
	case "azurerm_private_dns_mx_record":
		vals, err := decoderPrivateDnsMXRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDNSMXRecord(vals).component()
	case "azurerm_private_dns_ptr_record":
		vals, err := decoderPrivateDnsPTRRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDNSPTRRecord(vals).component()
	case "azurerm_private_dns_srv_record":
		vals, err := decoderPrivateDnsSRVRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDNSSRVRecord(vals).component()
	case "azurerm_private_dns_txt_record":
		vals, err := decoderPrivateDnsTXTRecord(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDNSTXTRecord(vals).component()
	case "azurerm_private_dns_zone":
		vals, err := decoderPrivateDnsZone(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPrivateDNSZone(vals).component()
	case "azurerm_cosmosdb_table":
		vals, err := decodeCosmosdbTableValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbTable(vals).Components()
	case "azurerm_cosmosdb_sql_database":
		vals, err := decodeCosmosdbSqlDatabaseValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbSqlDatabase(vals).Components()
	case "azurerm_cosmosdb_sql_container":
		vals, err := decodeCosmosdbSqlContainerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbSqlContainer(vals).Components()
	case "azurerm_cosmosdb_gremlin_database":
		vals, err := decodeCosmosdbGremlinDatabaseValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbGremlinDatabase(vals).Components()
	case "azurerm_cosmosdb_gremlin_graph":
		vals, err := decodeCosmosdbGremlinGraphValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbGremlinGraph(vals).Components()
	case "azurerm_cosmosdb_mongo_database":
		vals, err := decodeCosmosdbMongoDatabaseValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbMongoDatabase(vals).Components()
	case "azurerm_cosmosdb_cassandra_keyspace":
		vals, err := decodeCosmosdbCassandraKeyspaceValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbCassandraKeyspace(vals).Components()
	case "azurerm_cosmosdb_cassandra_table":
		vals, err := decodeCosmosdbCassandraTableValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbCassandraTable(vals).Components()
	case "azurerm_cosmosdb_mongo_collection":
		vals, err := decodeCosmosdbMongoCollectionValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newCosmosdbMongoCollection(vals).Components()
	case "azurerm_mariadb_server":
		vals, err := decodeMariadbServerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newMariadbServer(vals).Components()
	case "azurerm_sql_database":
		vals, err := decodeSqlDatabaseValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newSQLDatabase(vals).Components()
	case "azurerm_mssql_database":
		vals, err := decodeMssqlDatabaseValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newMssqlDatabase(vals).Components()
	case "azurerm_sql_managed_instance":
		vals, err := decodeSqlManagedInstanceValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newSqlManagedInstance(vals).Components()
	case "azurerm_mssql_managed_instance":
		vals, err := decodeMssqlManagedInstanceValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newMssqlManagedInstance(vals).Components()
	case "azurerm_mysql_server":
		vals, err := decodeMysqlServerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newMysqlServer(vals).Components()
	case "azurerm_postgresql_server":
		vals, err := decodePostgresqlServerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPostgresqlServer(vals).Components()
	case "azurerm_postgresql_flexible_server":
		vals, err := decodePostgresqlFlexibleServerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newPostgresqlFlexibleServer(vals).Components()
	case "azurerm_mysql_flexible_server":
		vals, err := decodeMysqlFlexibleServerValues(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newMysqlFlexibleServer(vals).Components()
	case "azurerm_kubernetes_cluster":
		vals, err := decoderKubernetesCluster(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.NewAzureRMKubernetesCluster(vals).Components()
	case "azurerm_kubernetes_cluster_node_pool":
		vals, err := decoderKubernetesClusterNodePool(tfRes.Values)
		if err != nil {
			p.logger.Debug(fmt.Sprintf("%v.%v could not be decoded: %s", tfRes.Type, tfRes.Name, err.Error()))
			return nil
		}
		return p.newAzureRMKubernetesClusterNodePool(vals).Components()
	default:
		return nil
	}
}

// getLocationName will return the location name from the location display name (ex: UK West -> ukwest)
// if the l is not found it'll return the l again meaning is not found or already a name
func getLocationName(l string) string {
	ln, ok := locationDisplayToName[l]
	if !ok {
		return l
	}
	return ln
}

func GetCostComponentNamesAndSetLogger(costComponents []resource.Component, logger *zap.Logger) {
	costComponentsName := ""
	for k, v := range costComponents {
		if k != 0 {
			costComponentsName += ","
		}
		costComponentsName += v.Name
	}
	logger.Info(fmt.Sprintf("cost components name : %v", costComponentsName))
}

func regionToZone(region string) string {
	return map[string]string{
		"westus":             "Zone 1",
		"westus2":            "Zone 1",
		"eastus":             "Zone 1",
		"centralus":          "Zone 1",
		"centraluseuap":      "Zone 1",
		"southcentralus":     "Zone 1",
		"northcentralus":     "Zone 1",
		"westcentralus":      "Zone 1",
		"eastus2":            "Zone 1",
		"eastus2euap":        "Zone 1",
		"brazilsouth":        "Zone 3",
		"brazilus":           "Zone 3",
		"northeurope":        "Zone 1",
		"westeurope":         "Zone 1",
		"eastasia":           "Zone 2",
		"southeastasia":      "Zone 2",
		"japanwest":          "Zone 2",
		"japaneast":          "Zone 2",
		"koreacentral":       "Zone 2",
		"koreasouth":         "Zone 2",
		"southindia":         "Zone 5",
		"westindia":          "Zone 5",
		"centralindia":       "Zone 5",
		"australiaeast":      "Zone 4",
		"australiasoutheast": "Zone 4",
		"canadacentral":      "Zone 1",
		"canadaeast":         "Zone 1",
		"uksouth":            "Zone 1",
		"ukwest":             "Zone 1",
		"francecentral":      "Zone 1",
		"francesouth":        "Zone 1",
		"australiacentral":   "Zone 4",
		"australiacentral2":  "Zone 4",
		"uaecentral":         "Zone 1",
		"uaenorth":           "Zone 1",
		"southafricanorth":   "Zone 1",
		"southafricawest":    "Zone 1",
		"switzerlandnorth":   "Zone 1",
		"switzerlandwest":    "Zone 1",
		"germanynorth":       "Zone 1",
		"germanywestcentral": "Zone 1",
		"norwayeast":         "Zone 1",
		"norwaywest":         "Zone 1",
		"brazilsoutheast":    "Zone 3",
		"westus3":            "Zone 1",
		"eastusslv":          "Zone 1",
		"swedencentral":      "Zone 1",
		"swedensouth":        "Zone 1",
	}[region]
}
