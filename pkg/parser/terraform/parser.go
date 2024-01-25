package terraform

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/pkg/usage"
	"io"
	"reflect"
	"strings"
)

// Plan is a representation of a Terraform plan file.
type Plan struct {
	providerInitializers map[string]ProviderInitializer
	usage                usage.Usage

	Configuration Configuration       `json:"configuration"`
	PriorState    *State              `json:"prior_state"`
	PlannedValues Values              `json:"planned_values"`
	Variables     map[string]Variable `json:"variables"`
}

// SetUsage will set the usage of the plan
func (p *Plan) SetUsage(u usage.Usage) { p.usage = u }

// NewPlan returns an empty Plan.
func NewPlan(providerInitializers ...ProviderInitializer) *Plan {
	piMap := make(map[string]ProviderInitializer)
	for _, pi := range providerInitializers {
		for _, name := range pi.MatchNames {
			piMap[name] = pi
		}
	}
	plan := &Plan{providerInitializers: piMap}
	return plan
}

// Read reads the Plan file from the provider io.Reader.
func (p *Plan) Read(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(p); err != nil {
		return err
	}
	return nil
}

// ExtractPlannedQueries extracts a query.Resource slice from the `planned_values` part of the Plan.
func (p *Plan) ExtractPlannedQueries() ([]Resource, error) {
	providers, err := p.extractProviders()
	if err != nil {
		return nil, fmt.Errorf("unable to extract planned queries: %w", err)
	}
	q, err := p.extractResources(p.PlannedValues, providers)
	if err != nil {
		return nil, fmt.Errorf("failed to extract queries: %w", err)
	}
	return q, nil
}

// ExtractPriorQueries extracts a query.Resource slice from the `prior_state` part of the Plan.
func (p *Plan) ExtractPriorQueries() ([]Resource, error) {
	if p.PriorState == nil {
		return []Resource{}, nil
	}
	providers, err := p.extractProviders()
	if err != nil {
		return nil, fmt.Errorf("unable to extract prior queries: %w", err)
	}

	q, err := p.extractResources(p.PriorState.Values, providers)
	if err != nil {
		return nil, fmt.Errorf("failed to extract queries: %w", err)
	}

	return q, nil
}

// extractProviders returns a slice of initialized Provider instances that were found in plan's configuration.
func (p *Plan) extractProviders() (map[string]Provider, error) {
	providers := make(map[string]Provider)
	for name, provConfig := range p.Configuration.ProviderConfig {
		if pi, ok := p.providerInitializers[provConfig.Name]; ok {
			values, err := p.evaluateProviderConfigExpressions(provConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to read config of provider %q: %w", name, err)
			}
			prov, err := pi.Provider(values)
			if err != nil {
				return nil, err
			}
			if prov != nil {
				providers[name] = prov
			}
		}
	}
	if len(providers) == 0 {
		return nil, ErrNoProviders
	}
	return providers, nil
}

// extractResources iterates over every resource and passes each to the corresponding Provider to get the components.
// These are used to form a slice of resource queries that are then returned back to the caller.
func (p *Plan) extractResources(values Values, providers map[string]Provider) ([]Resource, error) {
	// Create a map to associate each resource with a Provider that
	// should be used to estimate it.
	resourceProviders := make(map[string]providerWithResourceValues)
	err := p.extractModuleConfiguration("", &p.Configuration.RootModule, providers, resourceProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to extract module (%s) configuraiotn: %w", "root_module", err)
	}

	resourcesMap := p.extractModuleResources(&values.RootModule, resourceProviders)

	var retry int
	var oldResourcesMap map[string][]Resource
	for retry < 50 {
		resourcesMap = p.extractReferences(resourcesMap)
		if mapsEqual(oldResourcesMap, resourcesMap) {
			break
		}
		oldResourcesMap = resourcesMap
		retry++
	}

	var resources []Resource
	for _, rss := range resourcesMap {
		resources = append(resources, rss...)
	}
	return resources, nil
}

type providerWithResourceValues struct {
	Provider Provider
	Values   map[string]interface{}
}

func (p *Plan) extractReferences(resourcesMap map[string][]Resource) map[string][]Resource {
	for _, resources := range resourcesMap {
		for i, res := range resources {
			res.Values["id"] = fmt.Sprintf("%s.id", res.Address)
			for key, val := range res.Values {
				if value, ok := val.(string); ok {
					ref := strings.Split(value, ".")
					if ref[0] == "*ref*" && len(ref) > 2 {
						if _, ok := resourcesMap[strings.Join(ref[1:len(ref)-1], ".")]; ok {
							if refValue, ok := resourcesMap[strings.Join(ref[1:len(ref)-1], ".")][0].Values[ref[len(ref)-1]]; ok {
								res.Values[key] = refValue
							}
						} else {
							res.Values[key] = nil
						}
					} else if ref[0] == "*each*" && len(ref) > 2 {
						if _, ok := resourcesMap[strings.Join(ref[1:len(ref)-1], ".")]; ok && len(resourcesMap[strings.Join(ref[1:len(ref)-1], ".")]) > i {
							if refValue, ok := resourcesMap[strings.Join(ref[1:len(ref)-1], ".")][i].Values[ref[len(ref)-1]]; ok {
								res.Values[key] = refValue
							}
						} else {
							res.Values[key] = nil
						}
					}
				} else {
					continue
				}
			}
		}
	}
	return resourcesMap
}

// extractModuleConfiguration iterates over all the modules included in the plan's configuration block and
// extracts the provider that should be used for each resource. This function calls itself recursively until
// data from the entire module tree is extracted. It takes the following arguments:
//   - prefix - the current module's address. Empty string signifies the root module.
//   - module - the module's configuration block itself.
//   - providers - map of provider name to Provider.
//   - resourceProviders - used as an output of this function, it's a map of resource addresses to their assigned
//     Provider and the values on the resource. This map should be passed empty and not nil.
func (p *Plan) extractModuleConfiguration(prefix string, module *ConfigurationModule, providers map[string]Provider, resourceProviders map[string]providerWithResourceValues) error {
	for _, res := range module.Resources {
		key := res.ProviderConfigKey
		if strings.Contains(key, ":") {
			parts := strings.Split(key, ":")
			key = parts[len(parts)-1]
		}

		addr := res.Address
		if prefix != "" {
			addr = fmt.Sprintf("module.%s.%s", prefix, addr)
		}

		if prov, ok := providers[key]; ok {
			resPrefix := fmt.Sprintf("module.%s", prefix)
			rv, err := p.evaluateResourceExpressions(resPrefix, res.ForEachExpression, res.Expressions, module.Variables)
			if err != nil {
				return fmt.Errorf("failed to evaluate resource expresions: %w", err)
			}
			resourceProviders[addr] = providerWithResourceValues{
				Provider: prov,
				Values:   rv,
			}
		}
	}

	for k, child := range module.ModuleCalls {
		if child.Module != nil {
			nextPrefix := k
			if prefix != "" {
				nextPrefix = fmt.Sprintf("%s.%s", prefix, k)
			}
			err := p.extractModuleConfiguration(nextPrefix, child.Module, providers, resourceProviders)
			if err != nil {
				return fmt.Errorf("failed to extract child (%s) module configuration: %w", nextPrefix, err)
			}
		}
	}
	return nil
}

// extractModuleResources iterates recursively over all the module's (and its descendants) resources. It uses the
// resourceProviders map to retrieve the correct Provider based on the resource address.
func (p *Plan) extractModuleResources(module *Module, resourceProviders map[string]providerWithResourceValues) map[string][]Resource {
	rss := make(map[string]Resource)
	resources := make(map[string][]Resource)
	for _, tfres := range module.Resources {
		pwrv, ok := resourceProviders[fmt.Sprintf("%s.%s", tfres.Type, tfres.Name)]
		if !ok {
			pwrv, ok = resourceProviders[tfres.Address]
			if !ok {
				continue
			}
		}
		for k, v := range pwrv.Values {
			if v == nil {
				continue
			}

			vv, ok := tfres.Values[k]
			if !ok {
				tfres.Values[k] = v
				continue
			}

			switch tv := v.(type) {
			case map[string]interface{}:
				for tk, ntv := range tv {
					if ntv != nil {
						vv.(map[string]interface{})[tk] = ntv
					}
				}
			case []interface{}:
				for i, iv := range tv {
					miv, ok := iv.(map[string]interface{})
					if !ok {
						tfres.Values[k] = v
						continue
					}
					// We'll assume if they have a nil value that the other
					// one is correct
					nmap := tfres.Values[k].([]interface{})[i].(map[string]interface{})
					for ntk, ntv := range miv {
						if ntv != nil {
							// We only set the new value to the map if it's not present
							// on the original one
							if _, ok := nmap[ntk]; !ok {
								nmap[ntk] = ntv
							}
						}
					}
					tfres.Values[k].([]interface{})[i] = nmap
				}
			default:
				// If it's a base type and we ignored the empty ones we'll
				// always take the resource value over the Provider value
				// as the resource could be the Prior hence we would need
				// use the old value and the Provider always have the new value
				continue
			}
		}
		rss[tfres.Address] = tfres
		tfres.Values[usage.Key] = p.usage.GetUsage(tfres.Type, tfres.Address)
	}

	for _, rs := range rss {
		// We know it's present as it has passed the previous loop
		name := fmt.Sprintf("%s.%s", rs.Type, rs.Name)
		if _, ok := resources[name]; ok {
			resources[name] = append(resources[name], rs)
		} else {
			resources[name] = []Resource{rs}
		}
	}
	for _, child := range module.ChildModules {
		childResources := p.extractModuleResources(child, resourceProviders)
		for name, rs := range childResources {
			if _, ok := resources[name]; ok {
				resources[name] = append(resources[name], rs...)
			} else {
				resources[name] = []Resource{}
				resources[name] = append(resources[name], rs...)
			}
		}
	}

	return resources
}

// evaluateProviderConfigExpressions returns evaluated values of provider's configuration block, whether a constant
// value or reference to a variable.
func (p *Plan) evaluateProviderConfigExpressions(config ProviderConfig) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	for name, e := range config.Expressions {
		if e.ConstantValue != "" {
			values[name] = e.ConstantValue
			continue
		}

		if len(e.References) < 1 {
			return nil, fmt.Errorf("config expression contains invalid reference")
		}

		ref := strings.Split(e.References[0], ".")
		if len(ref) < 2 {
			return nil, fmt.Errorf("config expression contains invalid reference")
		}

		varName := ref[1]
		v, ok := p.Variables[varName]
		if !ok || v.Value == "" {
			return nil, fmt.Errorf("required variable %q is not defined", varName)
		}
		values[name] = v.Value
	}
	return values, nil
}

// evaluateResourceExpressions returns evaluated values of resource's configuration block, whether a constant
// value or reference to a variable.
func (p *Plan) evaluateResourceExpressions(prefix string, forEach map[string]interface{}, config map[string]interface{}, variables map[string]Variable) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	for name, ex := range config {
		m, ok := ex.(map[string]interface{})
		if !ok {
			a, ok := ex.([]interface{})
			if !ok {
				// If it's not a map or [] then we just ignore it
				continue
			}
			values[name] = make([]interface{}, 0, 0)
			for _, c := range a {
				mc, ok := c.(map[string]interface{})
				if !ok {
					// It should always be a []map so if it's not then we ignore it
					// this should be values and if it's an array it means it's a block
					// that can be defined multiple times so it should always be map[]
					continue
				}
				av, err := p.evaluateResourceExpressions(prefix, forEach, mc, variables)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluateResourceExpressions on array: %w", err)
				}
				values[name] = append(values[name].([]interface{}), av)
			}
		}
		refs, ok := m["references"].([]interface{})
		if !ok {
			refs = make([]interface{}, 0, 0)
		}

		if len(m) > 0 && m["constant_value"] == nil && len(refs) == 0 {
			// Right now the only key identified empty has been `timeout`
			continue
		}
		if m["constant_value"] != nil {
			values[name] = m["constant_value"]
			continue
		}

		if len(refs) < 1 {
			continue
		}

		ref := strings.Split(refs[0].(string), ".")
		if len(ref) < 2 {
			return nil, fmt.Errorf("refernce %q has invalid format", refs[0])
		}
		if ref[0] == "each" {
			values[name] = fmt.Sprintf("*each*.%s.%s", forEach["references"].([]interface{})[0], ref[2])
			continue
		}
		// "local" variables are not set on the plan
		// so we ignore them
		if ref[0] == "local" {
			continue
		}

		// For now we do not want external module references or data references
		if ref[0] == "module" || ref[0] == "data" {
			continue
		}

		if ref[0] == "var" {
			varName := ref[1]
			v, ok := variables[varName]
			if !ok || v.Value == "" {
				return nil, fmt.Errorf("required variable %q is not defined", varName)
			}
			values[name] = v.Value
			continue
		}
		values[name] = fmt.Sprintf("%s", refs[0])
	}
	return values, nil
}

// mapsEqual checks if two maps are equal
func mapsEqual(resourcesMap1, resourcesMap2 map[string][]Resource) bool {
	if len(resourcesMap1) != len(resourcesMap2) {
		return false
	}

	for key, value1 := range resourcesMap1 {
		if value2, ok := resourcesMap2[key]; ok {
			if len(value1) != len(value2) {
				return false
			}
			for i, _ := range value1 {
				map1 := value1[i].Values
				map2 := value2[i].Values
				for k, v1 := range map1 {
					if v2, ok := map2[k]; ok {
						if !reflect.DeepEqual(v1, v2) {
							return false
						}
					} else {
						return false
					}
				}
			}
		} else {
			return false
		}
	}

	return true
}
