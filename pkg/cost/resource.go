package cost

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kaytu-io/pennywise/pkg/schema"
)

// Resource represents costs of a single cloud resource. Each Resource includes a Component map, keyed
// by the label.
type Resource struct {
	Provider    schema.ProviderName
	Type        string
	Components  map[string][]Component
	Skipped     bool
	IsSupported bool
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

// CostRows returns rows for resource components
// containing the components costs and total cost for the resource
func (re Resource) CostRows() ([]table.Row, error) {
	var rows []table.Row

	for _, comps := range re.Components {
		for _, c := range comps {
			if c.Cost().Decimal.Round(3).IntPart() == 0 {
				continue
			}
			var row table.Row
			row = append(row, faint.Sprint("└─ ")+c.Name, c.Rate.Decimal, c.HourlyQuantity, c.MonthlyQuantity, c.Unit, c.Cost().Decimal.Round(2))
			rows = append(rows, row)
		}
	}
	return rows, nil
}

// ResourceDiff is the difference in costs between prior and planned Resource. It contains a ComponentDiff
// map, keyed by the label.
type ResourceDiff struct {
	Address        string
	Provider       schema.ProviderName
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

// Valid returns true if there are no errors in all of the ResourceDiff components.
func (rd ResourceDiff) Valid() bool {
	for _, cd := range rd.ComponentDiffs {
		if (cd.Prior != nil && cd.Prior.Error != nil) || (cd.Planned != nil && cd.Planned.Error != nil) {
			return false
		}
	}
	return true
}
