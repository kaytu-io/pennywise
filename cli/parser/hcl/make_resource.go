package hcl

import (
	"fmt"
)

var makeResourceProcesses = map[string]MakeResourceProcess{
	"azurerm_snapshot": {
		Refs: []string{"source_uri"},
	},
}

type ResourceFunction func(Resource) (Resource, error)

type MakeResourceProcess struct {
	Refs      []string
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
			if key == ref {
				if refId == nil {
					break
				}
				res, err := findResource(rss, refId.(string))
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

func findResource(rss []Resource, id string) (*Resource, error) {
	for _, res := range rss {
		for key, value := range res.Values {
			if key == "id" || key == "self_link" {
				if value.(string) == id {
					return &res, nil
				} else {
					continue
				}
			}
		}
		return nil, fmt.Errorf("id field not found")
	}
	return nil, fmt.Errorf("resource not found")
}
