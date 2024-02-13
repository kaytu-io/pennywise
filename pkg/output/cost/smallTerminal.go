package cost

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaytu-io/pennywise/pkg/cost"
)

type SmallTerminalModel struct {
	state       *cost.ModularState
	parentModel *ResourcesModel
	label       string
	wSize       int
}

func (m SmallTerminalModel) Init() tea.Cmd { return nil }

func (m SmallTerminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			model, err := getResourcesModel(m.label, m.state, m.wSize, m.parentModel)
			if err != nil {
				panic(err)
			}
			return model, cmd
		}
	}
	return m, cmd
}

func (m SmallTerminalModel) View() string {
	return "Can't output with your terminal size\nShow anyway [press ENTER], you can also view in classic mode using --classic tag\n" +
		"Exit by pressing [ESC], q or [CTRL+C]"
}

func getSmallTerminalModelModel(label string, state *cost.ModularState, wSize int, parentModel *ResourcesModel) (tea.Model, error) {

	m := SmallTerminalModel{state, parentModel, label, wSize}
	return m, nil
}
