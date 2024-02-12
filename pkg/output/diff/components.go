package diff

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/leekchan/accounting"
	"golang.org/x/crypto/ssh/terminal"
)

type ComponentsModel struct {
	label          string
	table          table.Model
	resourcesModel ResourcesModel
}

func (m ComponentsModel) Init() tea.Cmd { return nil }

func (m ComponentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left":
			return m.resourcesModel, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ComponentsModel) View() string {
	output := "Navigate to resources by pressing ← Quit by pressing Q or [CTRL+C]\n\n"
	output += bold.Sprint(m.label) + "\n" + baseStyle.Render(m.table.View()) + "\n"
	output += "To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md"
	return output
}

func getComponentsModel(resourceName, resourceCost string, components map[string][]schema.ComponentDiff, resModel ResourcesModel) (tea.Model, error) {
	var longestName int
	for _, comps := range components {
		for _, c := range comps {
			if len(c.Component.Name) > longestName {
				longestName = len(c.Component.Name)
			}
		}
	}
	w, _, err := terminal.GetSize(0)
	if err != nil {
		return nil, err
	}
	if (longestName + 83) > w {
		longestName = w - 83
	}
	columns := []table.Column{
		{Title: "Name", Width: longestName + 11},
		{Title: "Unit Price", Width: 10},
		{Title: "Hourly Qty", Width: 20},
		{Title: "Monthly Qty", Width: 20},
		{Title: "Unit", Width: 10},
		{Title: "Monthly Cost", Width: 25},
	}

	var rows []table.Row
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	for _, comps := range components {
		for _, c := range comps {
			var row table.Row
			var componentName string
			var rateString string
			var hourlyCost string
			var monthlyCost string
			var costDiff string
			switch c.Action {
			case schema.ActionCreate:
				componentName = green.Sprint("+ ") + c.Component.Name
				rateString = c.Component.Rate.Decimal.String()
				hourlyCost = c.Component.HourlyQuantity.String()
				monthlyCost = c.Component.MonthlyQuantity.String()
				costDiff = "+" + ac.FormatMoney(c.CostDiff)
			case schema.ActionModify:
				componentName = yellow.Sprint("~ ") + c.Component.Name
				if !c.Component.Rate.Decimal.IsZero() {
					rateString = c.Component.Rate.Decimal.String() +
						fmt.Sprintf(" (%s -> %s)", c.CompareTo.Rate.Decimal.String(), c.Current.Rate.Decimal.String())
				}
				hourlyCost = c.Component.HourlyQuantity.String()
				if !c.Component.HourlyQuantity.IsZero() {
					hourlyCost = hourlyCost +
						fmt.Sprintf(" (%s -> %s)", c.CompareTo.HourlyQuantity.String(), c.Current.HourlyQuantity.String())
				}
				monthlyCost = c.Component.MonthlyQuantity.String()
				if !c.Component.MonthlyQuantity.IsZero() {
					monthlyCost = monthlyCost +
						fmt.Sprintf(" (%s -> %s)", c.CompareTo.MonthlyQuantity.String(), c.Current.MonthlyQuantity.String())
				}
				costDiff = ac.FormatMoney(c.CostDiff) +
					fmt.Sprintf(" (%s -> %s)", ac.FormatMoney(c.CompareTo.Cost().Decimal), ac.FormatMoney(c.Current.Cost().Decimal))
				if c.CostDiff.InexactFloat64() > 0 {
					costDiff = "+" + costDiff
				}
			case schema.ActionRemove:
				componentName = red.Sprint("- ") + c.Component.Name
				rateString = c.Component.Rate.Decimal.String()
				hourlyCost = c.Component.HourlyQuantity.String()
				monthlyCost = c.Component.MonthlyQuantity.String()
				costDiff = "-" + ac.FormatMoney(c.CostDiff)
			}
			row = append(row, componentName, rateString, hourlyCost,
				monthlyCost, c.Component.Unit, costDiff)
			rows = append(rows, row)
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#808080")).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderLeft(true).BorderBottom(false).BorderRight(false).BorderTop(false)
	t.SetStyles(s)
	m := ComponentsModel{fmt.Sprintf("%s, Resource Total Cost: %s", resourceName, resourceCost), t, resModel}
	return m, nil
}
