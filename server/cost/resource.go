package cost

import "fmt"

// Resource represents costs of a single cloud resource. Each Resource includes a Component map, keyed
// by the label.
type Resource struct {
	Provider   string
	Type       string
	Components map[string][]Component
	Skipped    bool
}

// Cost returns the sum of costs of every Component of this Resource.
// Error is returned if there is a mismatch in Component currency.
func (re Resource) Cost() (Cost, error) {
	var total Cost
	var err error
	for name, comp := range re.Components {
		for _, c := range comp {
			total, err = total.Add(c.Cost())
			if err != nil {
				return Zero, fmt.Errorf("failed to add cost of component %s: %w", name, err)
			}

		}
	}
	return total, nil
}

func (re Resource) CostString() (string, error) {
	cost, err := re.Cost()
	if err != nil {
		return "", err
	}
	costString := fmt.Sprintf("---- Total Resource Cost: %v", cost)
	for _, comps := range re.Components {
		for _, c := range comps {
			costString = fmt.Sprintf("%s\n-------- %s : %s", costString, c.Name, c.CostString())
		}
	}
	return costString, nil
}

// ResourceDiff is the difference in costs between prior and planned Resource. It contains a ComponentDiff
// map, keyed by the label.
type ResourceDiff struct {
	Address        string
	Provider       string
	Type           string
	ComponentDiffs map[string]*ComponentDiff
}

// Errors returns a map of Component errors keyed by the Component label.
func (rd ResourceDiff) Errors() map[string]error {
	errs := make(map[string]error)
	for label, cd := range rd.ComponentDiffs {
		if cd.Prior != nil && cd.Prior.Error != nil {
			errs[label] = cd.Prior.Error
		} else if cd.Planned != nil && cd.Planned.Error != nil {
			errs[label] = cd.Planned.Error
		}
	}
	return errs
}

// Valid returns true if there are no errors in all of the ResourceDiff's components.
func (rd ResourceDiff) Valid() bool {
	for _, cd := range rd.ComponentDiffs {
		if (cd.Prior != nil && cd.Prior.Error != nil) || (cd.Planned != nil && cd.Planned.Error != nil) {
			return false
		}
	}
	return true
}
