package cost

import (
	"fmt"
	"github.com/kaytu-io/pennywise/pkg/schema"
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

	return Cost{Currency: total.Currency, Decimal: total.Decimal.Round(3)}, nil
}

func (s *State) GetCostComponents() []Component {
	var components []Component
	for _, res := range s.Resources {
		for _, comp := range res.Components {
			for _, c := range comp {
				components = append(components, c.GetRounded())
			}
		}
	}
	return components
}

// CostString returns a string to show the breakdown of the costs for a state
// containing the resources and their components costs and total cost for the resources and the state
func (s *State) CostString() (string, error) {
	cost, err := s.Cost()
	if err != nil {
		return "", err
	}
	costString := fmt.Sprintf("- Total Cost (per month): %v", cost.Decimal.Round(3))
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

// EnsureResource creates Resource at the given address if it doesn't already exist.
func (s *State) EnsureResource(address, typ string, provider schema.ProviderName, skipped bool) {
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

// AddComponent adds the Component with given label to the Resource at given address.
func (s *State) AddComponent(resAddress, compLabel string, component Component) {
	if _, ok := s.Resources[resAddress].Components[compLabel]; !ok {
		s.Resources[resAddress].Components[compLabel] = []Component{}
	}
	s.Resources[resAddress].Components[compLabel] = append(s.Resources[resAddress].Components[compLabel], component)
}
