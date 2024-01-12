package main

import (
	"fmt"
	my_hcl "github.com/kaytu-io/pennywise/parser/my-hcl"
)

func main() {
	//cmd.Execute()
	tp := my_hcl.NewTerraformProject("./testdata/parser/snapshot")
	err := tp.FindFiles()
	if err != nil {
		panic(err)
	}
	err = tp.ParseProjectBlocks()
	if err != nil {
		panic(err)
	}
	fmt.Println(tp.MakeProjectMapStructure())
}
