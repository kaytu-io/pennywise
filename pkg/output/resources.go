package output

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaytu-io/pennywise/pkg/cost"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ResourcesModel struct {
	table     table.Model
	resources map[string]cost.Resource
}

func (m ResourcesModel) Init() tea.Cmd { return nil }

func (m ResourcesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			resource := m.resources[m.table.SelectedRow()[0]]
			compsModel, err := getComponentsModel(resource.Components, m)
			if err != nil {
				panic(err)
			}
			return compsModel, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ResourcesModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func getResourcesModel(resources map[string]cost.Resource) (tea.Model, error) {
	columns := []table.Column{
		{Title: "Name", Width: 165},
		{Title: "Monthly Cost", Width: 15},
	}

	var rows []table.Row

	for name, resource := range resources {
		cost, err := resource.Cost()
		if err != nil {
			return nil, err
		}
		rows = append(rows, []string{name, cost.Decimal.String()})
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

	m := ResourcesModel{t, resources}
	return m, nil
}
