package hcl

import (
	"fmt"
	"strings"
)

var makeResourceProcesses = map[string]MakeResourceProcess{
	"azurerm_snapshot": {
		Refs: []Reference{{RefValue: "source_uri", RefAttribute: "id"}},
	},
	"azurerm_lb_rule": {
		Refs: []Reference{{RefValue: "loadbalancer_id", RefAttribute: "id"}},
	},
	"azurerm_lb_outbound_rule": {
		Refs: []Reference{{RefValue: "loadbalancer_id", RefAttribute: "id"}},
	},
	"azurerm_storage_queue": {
		Refs: []Reference{{RefValue: "storage_account_name", RefAttribute: "azurerm_storage_account.name"}},
	},
	"azurerm_storage_share": {
		Refs: []Reference{{RefValue: "storage_account_name", RefAttribute: "azurerm_storage_account.name"}},
	},
	"azurerm_key_vault_key": {
		Refs: []Reference{{RefValue: "key_vault_id", RefAttribute: "id"}},
	},
	"azurerm_key_vault_certificate": {
		Refs: []Reference{{RefValue: "key_vault_id", RefAttribute: "id"}},
	},
	"azurerm_virtual_network_gateway_connection": {
		Refs: []Reference{{RefValue: "virtual_network_gateway_id", RefAttribute: "id"}},
	},
	"azurerm_virtual_network_peering": {
		Refs: []Reference{{RefValue: "virtual_network_name", RefAttribute: "name"}, {RefValue: "remote_virtual_network_id", RefAttribute: "id"}},
	},
	"azurerm_cdn_endpoint": {
		Refs: []Reference{{RefValue: "profile_name", RefAttribute: "name"}},
	},
	"azurerm_dns_a_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_aaaa_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_caa_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_cname_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_mx_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_ns_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_ptr_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_srv_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_txt_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_dns_zone": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_a_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_aaaa_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_cname_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_mx_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_ptr_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_srv_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_txt_record": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_private_dns_zone": {
		Refs: []Reference{{RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_table": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_sql_database": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_sql_container": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_gremlin_database": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_gremlin_graph": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_mongo_database": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_cassandra_keyspace": {
		Refs: []Reference{{RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_cosmosdb_cassandra_table": {
		Refs: []Reference{{RefValue: "cassandra_keyspace_id", RefAttribute: "azurerm_cosmosdb_cassandra_keyspace.id"}},
	},
	"azurerm_cosmosdb_mongo_collection": {
		Refs: []Reference{{RefValue: "database_name", RefAttribute: "azurerm_cosmosdb_mongo_database.name"}, {RefValue: "account_name", RefAttribute: "azurerm_cosmosdb_account.name"}, {RefValue: "resource_group_name", RefAttribute: "azurerm_resource_group.name"}},
	},
	"azurerm_mssql_database": {
		Refs: []Reference{{RefValue: "server_id", RefAttribute: "id"}},
	},
}

type ResourceFunction func(Resource) (Resource, error)

type Reference struct {
	RefValue     string
	RefAttribute string
}

type MakeResourceProcess struct {
	Refs      []Reference
	Functions map[string]ResourceFunction
}

func (p MakeResourceProcess) runFunctions(rs Resource) Resource {
	for name, f := range p.Functions {
		newRs, err := f(rs)
		if err != nil {
			fmt.Println(fmt.Sprintf("error while running function %s on resource%s: %s", name, rs.Address, err.Error()))
		} else {
			rs = newRs
		}
	}
	return rs
}

func (p MakeResourceProcess) setRefs(rss []Resource, rs Resource) Resource {
	for _, ref := range p.Refs {
		for key, refId := range rs.Values {
			if key == ref.RefValue {
				if refId == nil {
					break
				}
				if _, ok := refId.(string); !ok {
					break
				}
				res, err := findResource(rss, refId.(string), ref.RefAttribute)
				if err != nil {
					fmt.Println(fmt.Sprintf("error while setting ref %s for resource %s (ref=%s): %s", ref, rs.Address, refId, err.Error()))
				}
				rs.Values[key] = res
				break
			}
		}
	}
	return rs
}

func findResource(rss []Resource, id string, refAttribute string) (*Resource, error) {
	ref := strings.Split(refAttribute, ".")
	if len(ref) > 1 {
		refAttribute = ref[1]
	}
	for _, res := range rss {
		if len(ref) > 1 && !(ref[0] == res.Type) {
			continue
		}
		for key, value := range res.Values {
			if key == refAttribute {
				if value.(string) == id {
					return &res, nil
				} else {
					continue
				}
			}
		}
	}
	return nil, fmt.Errorf("resource not found")
}
