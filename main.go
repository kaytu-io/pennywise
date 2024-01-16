package main

import (
	"fmt"
	my_hcl "github.com/kaytu-io/pennywise/parser/my-hcl"
	"go.uber.org/zap"
)

func main() {
	//cmd.Execute()
	logger, err := zap.NewProduction()
	tp := my_hcl.NewTerraformProject("./testdata/parser/storage_queue", logger)
	err = tp.FindFiles()
	if err != nil {
		logger.Error(err.Error())
	}
	err = tp.ParseProjectBlocks()
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println(tp.MakeProjectMapStructure())
	fmt.Println("===========================")
	if diagsStr, ok := tp.Diags.Show(); ok {
		fmt.Println(diagsStr)
	}
}
