package cost

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FreeResourcesModel struct {
	table          table.Model
	resourcesModel ResourcesModel
}

func (m FreeResourcesModel) Init() tea.Cmd { return nil }

func (m FreeResourcesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m FreeResourcesModel) View() string {
	output := "Navigate to resources by pressing ← Quit by pressing Q or [CTRL+C]\n\n"
	output += bold.Sprint("Free Resources") + "\n" + baseStyle.Render(m.table.View()) + "\n"
	output += "To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md"
	return output
}

func getFreeResourcesModel(resModel ResourcesModel) (tea.Model, error) {
	columns := []table.Column{
		{Title: "Name", Width: resModel.longestName + 17},
	}

	var rows []table.Row

	for _, name := range resModel.freeResources {
		rows = append(rows, []string{name})
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

	m := FreeResourcesModel{t, resModel}
	return m, nil
}
