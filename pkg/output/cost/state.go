package cost

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaytu-io/pennywise/pkg/cost"
	"github.com/leekchan/accounting"
	"os"
)

func ShowStateCosts(s *cost.ModularState) error {
	totalCost, err := s.Cost()
	if err != nil {
		return err
	}
	var longestName int
	for name, _ := range s.Resources {
		if len(name) > longestName {
			longestName = len(name)
		}
	}
	for name, _ := range s.ChildModules {
		if len(name) > longestName {
			longestName = len(name)
		}
	}
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	label := fmt.Sprintf("Total Cost: %s", ac.FormatMoney(totalCost.Decimal))
	model, err := getResourcesModel(label, s, longestName, nil)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
