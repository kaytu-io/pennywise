package hcl

import (
	"fmt"
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
	"azurerm_virtual_network_gateway_connection": {
		Refs: []string{"virtual_network_gateway_id"},
	},
	"azurerm_virtual_network_peering": {
		Refs: []string{"virtual_network_name", "remote_virtual_network_id"},
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
				res, err := findResource(rss, refId.(string), ref.RefAttribute)
				if err != nil {
					fmt.Println(fmt.Sprintf("error while setting ref %s for resource %s: %s", ref, rs.Address, err.Error()))
				}
				rs.Values[key] = res
				break
			}
		}
	}
	return rs
}

func findResource(rss []Resource, id string, refAttribute string) (*Resource, error) {
	for _, res := range rss {
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
