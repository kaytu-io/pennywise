package main

import (
	"fmt"
	my_hcl "github.com/kaytu-io/pennywise/parser/my-hcl"
	"go.uber.org/zap"
)

func main() {
	//cmd.Execute()
	logger, err := zap.NewProduction()
	tp := my_hcl.NewTerraformProject("./testdata/parser/snapshot", logger)
	err = tp.FindFiles()
	if err != nil {
		logger.Error(err.Error())
	}
	err = tp.ParseProjectBlocks()
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println(tp.MakeProjectMapStructure())
}
