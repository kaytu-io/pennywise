package my_hcl

import (
	"fmt"
	"github.com/kaytu-io/pennywise-server/resource"
	usagePackage "github.com/kaytu-io/pennywise/usage"
)

// ParseHclResources parses a new terraform project and return provider name and resources
func ParseHclResources(path string, usage usagePackage.Usage) (resource.ProviderName, []Resource, error) {
	tp := newTerraformProject(path)
	err := tp.FindFiles()
	if err != nil {
		return "", nil, err
	}
	err = tp.FindProjectBlocks()
	if err != nil {
		return "", nil, err
	}
	mapStructure := tp.makeProjectMapStructure()
	provider, resources, err := extractResourcesFromMapStructure(mapStructure)
	if err != nil {
		return "", nil, err
	}
	if usage != nil {
		for _, res := range resources {
			res.addUsage(usage)
		}
	}
	fmt.Println(provider, resources)
	fmt.Println("===========================")
	if diagsStr, ok := tp.Diags.Show(); ok {
		fmt.Println(diagsStr)
	}
	return resource.ProviderName(provider), resources, nil
}
