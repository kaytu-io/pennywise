package cost

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UnsupportedModel struct {
	table          table.Model
	resourcesModel ResourcesModel
}

func (m UnsupportedModel) Init() tea.Cmd { return nil }

func (m UnsupportedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m UnsupportedModel) View() string {
	output := "Navigate to resources by pressing ← Quit by pressing Q or [CTRL+C]\n\n"
	output += bold.Sprint("Unsupported Resource Types") + "\n" + baseStyle.Render(m.table.View()) + "\n"
	output += "To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md"
	return output
}

func getUnsupportedModel(resModel ResourcesModel) (tea.Model, error) {
	columns := []table.Column{
		{Title: "Resource Type", Width: resModel.longestName},
		{Title: "Resources Count", Width: 15},
	}

	var rows []table.Row

	for name, resources := range resModel.unsupportedResources {
		rows = append(rows, []string{name, fmt.Sprintf("%d", len(resources))})
	}
	rows = sortRows(rows)

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

	m := UnsupportedModel{t, resModel}
	return m, nil
}
