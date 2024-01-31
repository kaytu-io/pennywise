package output

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FreeResourcesModel struct {
	viewport       viewport.Model
	table          table.Model
	resourcesModel ResourcesModel
}

func (m FreeResourcesModel) Init() tea.Cmd { return nil }

func (m FreeResourcesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "left", "q":
			return m.resourcesModel, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m FreeResourcesModel) View() string {
	return m.viewport.View() + "\n" + baseStyle.Render(m.table.View()) + "\n"
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
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	vp := viewport.New(30, 1)
	vp.SetContent(fmt.Sprintf("Free Resources"))
	m := FreeResourcesModel{vp, t, resModel}
	return m, nil
}
