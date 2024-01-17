package main

import (
	"github.com/kaytu-io/pennywise-server/resource"
	my_hcl "github.com/kaytu-io/pennywise/parser/my-hcl"
	"github.com/kaytu-io/pennywise/submission"
)

func main() {
	//cmd.Execute()
	provider, hclResources, err := my_hcl.ParseHclResources("./testdata/parser/storage_queue", nil)
	if err != nil {
		panic(err)
	}
	var resources []resource.ResourceDef
	for _, res := range hclResources {
		resources = append(resources, res.ToResourceDef(provider, nil))
	}
	sub, err := submission.CreateSubmission(resources)
	if err != nil {
		panic(err)
	}
	err = sub.StoreAsFile()
	if err != nil {
		panic(err)
	}
}
