package main

import (
	"fmt"
	my_hcl "github.com/kaytu-io/pennywise/parser/my-hcl"
)

func main() {
	//cmd.Execute()
	//tp := my_hcl.NewTerraformProject("./testdata/parser")
	tp := my_hcl.NewTerraformProject("../pennywise-server/testdata/azure/lb_rule")
	err := tp.FindFiles()
	if err != nil {
		panic(err)
	}
	err = tp.ParseProjectBlocks()
	if err != nil {
		panic(err)
	}
	for _, b := range tp.Blocks {
		fmt.Println("============================")
		fmt.Println(b)
		fmt.Println("---------------- attributes")
		b.ReadAttributes(tp.MappedBlocks)
	}
}
