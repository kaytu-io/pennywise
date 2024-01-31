package output

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaytu-io/pennywise/pkg/cost"
	"os"
)

func ShowStateCosts(s *cost.State) error {
	totalCost, err := s.Cost()
	if err != nil {
		return err
	}
	model, err := getResourcesModel(totalCost.Decimal.InexactFloat64(), s.Resources)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
