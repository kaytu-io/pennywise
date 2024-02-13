package diff

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/leekchan/accounting"
	"os"
)

var yellow = color.New(color.FgYellow)
var red = color.New(color.FgHiRed)
var green = color.New(color.FgHiGreen)

func ShowStateCosts(s *schema.ModularStateDiff) error {
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
	label := fmt.Sprintf("Total Diff: %s (%s -> %s)", ac.FormatMoney(s.NewCost.Sub(s.PriorCost)),
		ac.FormatMoney(s.PriorCost), ac.FormatMoney(s.NewCost))
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
