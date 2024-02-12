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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ResourcesModel struct {
	label                string
	table                table.Model
	resources            map[string]schema.ResourceDiff
	freeResources        []string
	unsupportedResources map[string][]string
	longestName          int
}

func (m ResourcesModel) Init() tea.Cmd { return nil }

func (m ResourcesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "enter":
			resourceName := m.table.SelectedRow()[0][11:]
			resource := m.resources[resourceName]
			compsModel, err := getComponentsModel(resourceName, m.table.SelectedRow()[1], resource.ComponentDiffs, m)
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
	output := "Navigate to details by pressing → or [ENTER] Quit by pressing Q or [CTRL+C]\n\n"
	output += bold.Sprint(m.label) + "\n" + baseStyle.Render(m.table.View()) + "\n"
	output += "To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md"
	return output
}

func getResourcesModel(label string, resources map[string]schema.ResourceDiff, longestName int) (tea.Model, error) {
	w, _, err := terminal.GetSize(0)
	if err != nil {
		return nil, err
	}
	if (longestName + 26) > w {
		return getSmallTerminalModelModel(label, resources, w-29)
	}
	columns := []table.Column{
		{Title: "Name", Width: longestName + 11},
		{Title: "Monthly Cost", Width: 30},
	}

	var rows []table.Row
	var freeResources []string
	unsupportedServices := make(map[string][]string)
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	for name, resource := range resources {
		if !resource.IsSupported && resource.Type != "" {
			if _, ok := unsupportedServices[resource.Type]; !ok {
				unsupportedServices[resource.Type] = []string{}
			}
			unsupportedServices[resource.Type] = append(unsupportedServices[resource.Type], name)
			continue
		}
		if resource.ComponentDiffs == nil {
			freeResources = append(freeResources, name)
			continue
		}
		var costDiff string
		switch resource.Action {
		case schema.ActionCreate:
			name = green.Sprint("+ ") + name
			costDiff = ac.FormatMoney(resource.NewCost)
		case schema.ActionModify:
			name = yellow.Sprint("~ ") + name
			costDiff = ac.FormatMoney(resource.NewCost.Sub(resource.PriorCost)) +
				fmt.Sprintf(" (%s -> %s)", ac.FormatMoney(resource.PriorCost), ac.FormatMoney(resource.NewCost))
			if resource.NewCost.Sub(resource.PriorCost).InexactFloat64() > 0 {
				costDiff = "+" + costDiff
			}
		case schema.ActionRemove:
			name = red.Sprint("- ") + name
			costDiff = ac.FormatMoney(resource.PriorCost)
		}
		rows = append(rows, []string{name, costDiff})
	}
	columns = append(columns, table.Column{Title: "", Width: 1})
	for i, _ := range rows {
		rows[i] = append(rows[i], "→")
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

	m := ResourcesModel{label, t, resources, freeResources, unsupportedServices, longestName}
	return m, nil
}
