package output

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaytu-io/pennywise/pkg/cost"
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
	output := "Navigate to resources by pressing ←  Quit by pressing Q or [CTRL+C]\n\n"
	output += bold.Sprint(m.label) + "\n" + baseStyle.Render(m.table.View()) + "\n"
	return output
}

func getComponentsModel(resourceName, resourceCost string, components map[string][]cost.Component, resModel ResourcesModel) (tea.Model, error) {
	var longestName int
	for _, comps := range components {
		for _, c := range comps {
			if len(c.Name) > longestName {
				longestName = len(c.Name)
			}
		}
	}
	w, _, err := terminal.GetSize(0)
	if err != nil {
		return nil, err
	}
	if (longestName + 70) > w {
		longestName = w - 70
	}
	columns := []table.Column{
		{Title: "Name", Width: longestName},
		{Title: "Unit Price", Width: 15},
		{Title: "Hourly Qty", Width: 15},
		{Title: "Monthly Qty", Width: 15},
		{Title: "Unit", Width: 10},
		{Title: "Monthly Cost", Width: 15},
	}

	var rows []table.Row

	for _, comps := range components {
		for _, c := range comps {
			var row table.Row
			row = append(row, c.Name, c.Rate.Decimal.String(), c.HourlyQuantity.String(), c.MonthlyQuantity.String(), c.Unit, c.Cost().Decimal.String())
			rows = append(rows, row)
		}
	}
	rows = sortRows(rows)
	rows = makeNumbersAccounting(rows)

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
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m := ComponentsModel{fmt.Sprintf("%s, Resource Total Cost: %s", resourceName, resourceCost), t, resModel}
	return m, nil
}
