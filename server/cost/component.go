package cost

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// Component describes the pricing of a single resource cost component. This includes Rate and Quantity
// and allows for final cost computation.
type Component struct {
	Name            string
	MonthlyQuantity decimal.Decimal
	HourlyQuantity  decimal.Decimal
	Unit            string
	Rate            Cost
	Details         []string
	Usage           bool

	Error error
}

func (c Component) GetRounded() Component {
	return Component{
		Name:            c.Name,
		MonthlyQuantity: c.MonthlyQuantity.Round(3),
		HourlyQuantity:  c.HourlyQuantity.Round(3),
		Unit:            c.Unit,
		Rate:            Cost{Decimal: c.Rate.Decimal.Round(3), Currency: c.Rate.Currency},
		Details:         c.Details,
		Usage:           c.Usage,

		Error: c.Error,
	}
}

// Cost returns the cost of this component (Rate multiplied by Quantity).
func (c Component) Cost() Cost {
	if !c.MonthlyQuantity.IsZero() {
		return c.Rate.MulDecimal(c.MonthlyQuantity)
	} else if !c.HourlyQuantity.IsZero() {
		return c.Rate.MulDecimal(c.HourlyQuantity.Mul(HoursPerMonth))
	} else {
		return Zero
	}
}

func (c Component) CostString() string {
	var str string
	str = fmt.Sprintf("%v", c.Cost().Decimal.Round(3))
	if !c.MonthlyQuantity.IsZero() {
		str = fmt.Sprintf("%s (%v monthly cost", str, c.Rate.Decimal.Round(3))
		if c.Unit != "" {
			str = fmt.Sprintf("%s per %s", str, c.Unit)
		}
		str = fmt.Sprintf("%s, qty: %s)", str, c.MonthlyQuantity.Round(3))
	} else if !c.HourlyQuantity.IsZero() {
		str = fmt.Sprintf("%s (%v hourly cost", str, c.Rate.Decimal.Round(3))
		if c.Unit != "" {
			str = fmt.Sprintf("%s per %s", str, c.Unit)
		}
		str = fmt.Sprintf("%s, qty: %s)", str, c.HourlyQuantity.Mul(HoursPerMonth).Round(3))
	} else {
		return fmt.Sprintf("No cost")
	}
	return str
}

// ComponentDiff is a difference between the Prior and Planned Component.
type ComponentDiff struct {
	Prior, Planned *Component
}

// Valid returns true if there are no errors in both the Planned and Prior components.
func (cd ComponentDiff) Valid() bool {
	return !((cd.Prior != nil && cd.Prior.Error != nil) || (cd.Planned != nil && cd.Planned.Error != nil))
}
