package diff

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"os"
)

func ShowStateCosts(s *schema.StateDiff) error {
	var longestName int
	for name, _ := range s.Resources {
		if len(name) > longestName {
			longestName = len(name)
		}
	}

	model, err := getResourcesModel(s.CostDiff.InexactFloat64(), s.Resources, longestName)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
