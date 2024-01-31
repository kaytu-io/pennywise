package output

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaytu-io/pennywise/pkg/cost"
)

type ComponentsModel struct {
	table          table.Model
	resourcesModel ResourcesModel
}

func (m ComponentsModel) Init() tea.Cmd { return nil }

func (m ComponentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m.resourcesModel, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ComponentsModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func getComponentsModel(components map[string][]cost.Component, resModel ResourcesModel) (tea.Model, error) {
	columns := []table.Column{
		{Title: "Name", Width: 50},
		{Title: "Unit Price", Width: 30},
		{Title: "Hourly Qty", Width: 25},
		{Title: "Monthly Qty", Width: 25},
		{Title: "Unit", Width: 25},
		{Title: "Monthly Cost", Width: 25},
	}

	var rows []table.Row

	for _, comps := range components {
		for _, c := range comps {
			if c.Cost().Decimal.Round(3).IntPart() == 0 {
				continue
			}
			var row table.Row
			row = append(row, c.Name, c.Rate.Decimal.String(), c.HourlyQuantity.String(), c.MonthlyQuantity.String(), c.Unit, c.Cost().Decimal.Round(2).String())
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
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := ComponentsModel{t, resModel}
	return m, nil
}
