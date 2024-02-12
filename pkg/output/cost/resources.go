package cost

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kaytu-io/pennywise/pkg/cost"
	"github.com/leekchan/accounting"
	"golang.org/x/crypto/ssh/terminal"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ResourcesModel struct {
	label                string
	table                table.Model
	state                *cost.ModularState
	parentModel          *ResourcesModel
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
		case "left":
			if m.parentModel != nil {
				return *m.parentModel, cmd
			}
		case "right", "enter":
			name := m.table.SelectedRow()[0]
			if name == "Free Resources" {
				freeResourcesModel, err := getFreeResourcesModel(m)
				if err != nil {
					panic(err)
				}
				return freeResourcesModel, cmd
			} else if name == "Unsupported" {
				unsupportedModel, err := getUnsupportedModel(m)
				if err != nil {
					panic(err)
				}
				return unsupportedModel, cmd
			}
			if resource, ok := m.state.Resources[name]; ok {
				compsModel, err := getComponentsModel(name, m.table.SelectedRow()[1], resource.Components, m)
				if err != nil {
					panic(err)
				}
				return compsModel, cmd
			} else {
				module := m.state.ChildModules[name]
				moduleCost, err := module.Cost()
				if err != nil {
					panic(err)
				}
				var longestName int
				for n, _ := range module.Resources {
					if len(n) > longestName {
						longestName = len(n)
					}
				}
				for n, _ := range module.ChildModules {
					if len(n) > longestName {
						longestName = len(n)
					}
				}
				resModel, err := getResourcesModel(moduleCost.Decimal.InexactFloat64(), &module, longestName, &m)
				if err != nil {
					panic(err)
				}
				return resModel, cmd
			}
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

func getResourcesModel(totalCost float64, state *cost.ModularState, longestName int, parentModel *ResourcesModel) (tea.Model, error) {
	w, _, err := terminal.GetSize(0)
	if err != nil {
		return nil, err
	}
	if (longestName + 33) > w {
		return getSmallTerminalModelModel(totalCost, state, w-36, parentModel)
	}
	columns := []table.Column{
		{Title: "Name", Width: longestName},
		{Title: "Resources", Width: 10},
		{Title: "Monthly Cost", Width: 12},
	}

	var rows []table.Row
	var freeResources []string
	unsupportedServices := make(map[string][]string)

	for name, module := range state.ChildModules {
		cost, err := module.Cost()
		if err != nil {
			return nil, err
		}
		rows = append(rows, []string{name, fmt.Sprintf("%d", module.TotalResourcesCount()), cost.Decimal.String()})
	}

	for name, resource := range state.Resources {
		if !resource.IsSupported && resource.Type != "" {
			if _, ok := unsupportedServices[resource.Type]; !ok {
				unsupportedServices[resource.Type] = []string{}
			}
			unsupportedServices[resource.Type] = append(unsupportedServices[resource.Type], name)
			continue
		}
		cost, err := resource.Cost()
		if err != nil {
			return nil, err
		}
		if resource.Components == nil {
			freeResources = append(freeResources, name)
			continue
		}
		rows = append(rows, []string{name, "", cost.Decimal.String()})
	}
	if len(freeResources) > 0 {
		rows = append(rows, []string{"Free Resources", fmt.Sprintf("%d", len(freeResources)), "0"})
	}
	if len(unsupportedServices) > 0 {
		rows = append(rows, []string{"Unsupported", fmt.Sprintf("%d", len(unsupportedServices)), "0"})
	}
	rows = sortRows(rows)
	rows = makeNumbersAccounting(rows)
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

	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	m := ResourcesModel{fmt.Sprintf("Total cost: %s", ac.FormatMoney(totalCost)), t, state, parentModel, freeResources, unsupportedServices, longestName}
	return m, nil
}
