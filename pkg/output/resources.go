package output

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
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
	viewport             viewport.Model
	table                table.Model
	resources            map[string]cost.Resource
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
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "enter":

			resourceName := m.table.SelectedRow()[0]
			if resourceName == "Free Resources" {
				freeResourcesModel, err := getFreeResourcesModel(m)
				if err != nil {
					panic(err)
				}
				return freeResourcesModel, cmd
			} else if resourceName == "Unsupported" {
				unsupportedModel, err := getUnsupportedModel(m)
				if err != nil {
					panic(err)
				}
				return unsupportedModel, cmd
			}
			resource := m.resources[resourceName]
			compsModel, err := getComponentsModel(resourceName, m.table.SelectedRow()[1], resource.Components, m)
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
	return m.viewport.View() + "\n" + baseStyle.Render(m.table.View()) + "\n"
}

func getResourcesModel(totalCost float64, resources map[string]cost.Resource, longestName int) (tea.Model, error) {
	w, _, err := terminal.GetSize(0)
	if err != nil {
		return nil, err
	}
	if (longestName + 20) > w {
		return getSmallTerminalModelModel(totalCost, resources, w-23)
	}
	columns := []table.Column{
		{Title: "Name", Width: longestName},
		{Title: "Monthly Cost", Width: 15},
	}

	var rows []table.Row
	var freeResources []string
	unsupportedServices := make(map[string][]string)

	for name, resource := range resources {
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
		if cost.Decimal.InexactFloat64() == 0 {
			freeResources = append(freeResources, name)
			continue
		}
		rows = append(rows, []string{name, cost.Decimal.String()})
	}
	if len(freeResources) > 0 {
		rows = append(rows, []string{"Free Resources", "0"})
	}
	if len(unsupportedServices) > 0 {
		rows = append(rows, []string{"Unsupported", "0"})
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

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	vp := viewport.New(30, 1)
	vp.SetContent(fmt.Sprintf("Total cost: %s", ac.FormatMoney(totalCost)))
	m := ResourcesModel{vp, t, resources, freeResources, unsupportedServices, longestName}
	return m, nil
}
