package cost

import (
	"context"
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/backend"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
)

// State represents a collection of all the Resource costs (either prior or planned.) It is not tied to any specific
// cloud provider or IaC tool. Instead, it is a representation of a snapshot of cloud resources at a given point
// in time, with their associated costs.
type State struct {
	Resources map[string]Resource
}

// Errors that might be returned from NewState if either a product or a price are not found.
var (
	ErrProductNotFound = fmt.Errorf("product not found")
	ErrPriceNotFound   = fmt.Errorf("price not found")
)

// NewState returns a new State from a query.Resource slice by using the Backend to fetch the pricing data.
func NewState(ctx context.Context, backend backend.Backend, resources []query.Resource) (*State, error) {
	state := &State{Resources: make(map[string]Resource)}
	if len(resources) == 0 {
		return nil, resource.ErrNoResources
	}
	for _, res := range resources {
		// Mark the Resource as skipped if there are no valid Components.
		state.ensureResource(res.Address, res.Provider, res.Type, len(res.Components) == 0)
		for _, comp := range res.Components {
			prods, err := backend.Products().Filter(ctx, comp.ProductFilter)
			if err != nil {
				state.addComponent(res.Address, comp.Name, Component{Error: err})
				continue
			}
			if len(prods) < 1 {
				fmt.Println("====================")
				fmt.Println("No product", comp.Name)
				fmt.Println("location", *comp.ProductFilter.Location)
				for _, attr := range comp.ProductFilter.AttributeFilters {
					if attr.Value != nil {
						fmt.Println(attr.Key, *attr.Value)
					}
					if attr.ValueRegex != nil {
						fmt.Println(attr.Key, *attr.ValueRegex)
					}
				}
				state.addComponent(res.Address, comp.Name, Component{Error: ErrProductNotFound})
				continue
			}
			prices, err := backend.Prices().Filter(ctx, prods, comp.PriceFilter)
			if err != nil {
				state.addComponent(res.Address, comp.Name, Component{Error: err})
				continue
			}
			if len(prices) < 1 {
				fmt.Println("NOT FOUND", comp.Name, comp.PriceFilter)
				fmt.Println("=====PRODS")
				for _, prod := range prods {
					fmt.Println(*prod)
				}
				state.addComponent(res.Address, comp.Name, Component{Error: ErrPriceNotFound})
				continue
			}

			component := Component{
				Name:            comp.Name,
				MonthlyQuantity: comp.MonthlyQuantity,
				HourlyQuantity:  comp.HourlyQuantity,
				Unit:            comp.Unit,
				Rate:            Cost{Decimal: prices[0].Value, Currency: prices[0].Currency},
				Details:         comp.Details,
				Usage:           comp.Usage,
			}

			state.addComponent(res.Address, comp.Name, component)
		}
	}

	return state, nil
}

// Cost returns the sum of the costs of every Resource included in this State.
// Error is returned if there is a mismatch in resource currencies.
func (s *State) Cost() (Cost, error) {
	var total Cost
	for name, re := range s.Resources {
		rCost, err := re.Cost()
		if err != nil {
			return Zero, fmt.Errorf("failed to get cost of resource %s: %w", name, err)
		}
		total, err = total.Add(rCost)
		if err != nil {
			return Zero, fmt.Errorf("failed to add cost of resource %s: %w", name, err)
		}
	}

	return total, nil
}

func (s *State) CostString() (string, error) {
	cost, err := s.Cost()
	if err != nil {
		return "", err
	}
	costString := fmt.Sprintf("- Total Cost (per month): %v", cost)
	for name, rs := range s.Resources {
		rsCostString, err := rs.CostString()
		if err != nil {
			return "", err
		}
		costString = fmt.Sprintf("%s\n--- Costs for %s :", costString, name)
		costString = fmt.Sprintf("%s\n%s", costString, rsCostString)

	}
	return costString, nil
}

// ensureResource creates Resource at the given address if it doesn't already exist.
func (s *State) ensureResource(address, provider, typ string, skipped bool) {
	if _, ok := s.Resources[address]; !ok {
		res := Resource{
			Provider: provider,
			Type:     typ,
			Skipped:  skipped,
		}

		if !skipped {
			res.Components = make(map[string][]Component)
		}

		s.Resources[address] = res
	}
}

// addComponent adds the Component with given label to the Resource at given address.
func (s *State) addComponent(resAddress, compLabel string, component Component) {
	if _, ok := s.Resources[resAddress].Components[compLabel]; !ok {
		s.Resources[resAddress].Components[compLabel] = []Component{}
	}
	s.Resources[resAddress].Components[compLabel] = append(s.Resources[resAddress].Components[compLabel], component)
}
