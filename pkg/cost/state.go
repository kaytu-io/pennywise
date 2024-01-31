package cost

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

var primary = color.New(color.FgHiCyan)

var yellow = color.New(color.FgYellow)
var red = color.New(color.FgHiRed)
var green = color.New(color.FgHiGreen)

var bold = color.New(color.Bold)
var faint = color.New(color.Faint)
var underline = color.New(color.Underline)

var primaryLink = color.New(color.Underline).Add(color.Bold)

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
	var costString string

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateRows = false
	t.Style().Options.SeparateHeader = false
	t.Style().Format.Header = text.FormatDefault

	var columns []table.ColumnConfig
	i := 1
	var headers table.Row
	headers = append(headers, underline.Sprint("Name"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignLeft,
		AlignHeader: text.AlignLeft,
	})
	i++

	headers = append(headers, underline.Sprint("Unit Price"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignLeft,
		AlignHeader: text.AlignLeft,
	})
	i++
	headers = append(headers, underline.Sprint("Hourly Qty"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignLeft,
		AlignHeader: text.AlignLeft,
	})
	i++
	headers = append(headers, underline.Sprint("Monthly Qty"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignLeft,
		AlignHeader: text.AlignLeft,
	})
	i++

	headers = append(headers, underline.Sprint("Unit"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignRight,
		AlignHeader: text.AlignRight,
	})
	i++

	headers = append(headers, underline.Sprint("Monthly Cost"))
	columns = append(columns, table.ColumnConfig{
		Number:      i,
		Align:       text.AlignRight,
		AlignHeader: text.AlignRight,
	})
	i++

	t.AppendRow(table.Row{""})

	t.SetColumnConfigs(columns)
	t.AppendHeader(headers)

	var unsupportedServices []string
	cost, err := s.Cost()
	if err != nil {
		return "", err
	}

	for name, rs := range s.Resources {
		if !rs.IsSupported {
			unsupportedServices = append(unsupportedServices, rs.Type)
			continue
		}
		cost, err := rs.Cost()
		if err != nil {
			return "", err
		}
		var row table.Row
		row = append(row, bold.Sprint(name), "", "", "", "", cost.Decimal.Round(2))
		costRows, err := rs.CostRows()
		if err != nil {
			return "", err
		}
		t.AppendRow(row)
		for _, r := range costRows {
			t.AppendRow(r)
		}
	}

	costString = t.Render()
	costString += "\n──────────────────────────────────\n"
	costString += fmt.Sprintf("%s:    %v", bold.Sprint("Total Cost (per month)"), cost.Decimal.Round(2))
	if len(unsupportedServices) == 3 {
		costString = fmt.Sprintf("%s\n- Resource types %s, %s and %s not supported", costString, unsupportedServices[0], unsupportedServices[1], unsupportedServices[2])
	} else if len(unsupportedServices) == 2 {
		costString = fmt.Sprintf("%s\n- Resource types %s and %s not supported", costString, unsupportedServices[0], unsupportedServices[1])
	} else if len(unsupportedServices) == 1 {
		costString = fmt.Sprintf("%s\n- Resource type %s not supported", costString, unsupportedServices[0])
	} else if len(unsupportedServices) > 3 {
		costString = fmt.Sprintf("%s\n- Resource types %s, %s, %s and %d other Resource types not supported", costString, unsupportedServices[0], unsupportedServices[1], unsupportedServices[2], len(unsupportedServices)-3)
	}

	return costString, nil
}

// EnsureResource creates Resource at the given address if it doesn't already exist.
func (s *State) EnsureResource(address, typ string, provider schema.ProviderName, skipped, isSupported bool) {
	if _, ok := s.Resources[address]; !ok {
		res := Resource{
			Provider:    provider,
			Type:        typ,
			Skipped:     skipped,
			IsSupported: isSupported,
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
