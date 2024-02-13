package cost

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Resource represents costs of a single cloud resource. Each Resource includes a Component map, keyed
// by the label.
type Resource struct {
	Address     string
	Provider    string
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
			var row table.Row
			row = append(row, faint.Sprint("└─ ")+c.Name, c.Rate.Decimal, c.HourlyQuantity, c.MonthlyQuantity, c.Unit, c.Cost().Decimal)
			rows = append(rows, row)
		}
	}
	return rows, nil
}
