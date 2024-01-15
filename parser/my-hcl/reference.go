package my_hcl

import (
	"fmt"
)

type Reference struct {
	blockType BlockType
	labels    []string
	key       string
}

type BlockType struct {
	name             string
	refName          string
	hasKey           bool
	getValueFunction func(block Block) (*Attribute, error)
}

var (
	BlockTypeResource = BlockType{
		name:    "resource",
		refName: "",
		hasKey:  true,
	}
	BlockTypeVariable = BlockType{
		name:             "variable",
		refName:          "var",
		hasKey:           false,
		getValueFunction: getVariableValue,
	}
	BlockTypeData = BlockType{
		name:    "data",
		refName: "data",
		hasKey:  false,
	}
	BlockTypeLocal = BlockType{
		name:    "locals",
		refName: "local",
		hasKey:  true,
	}
	BlockTypeProvider = BlockType{
		name:    "provider",
		refName: "provider",
		hasKey:  true,
	}
	BlockTypeOutput = BlockType{
		name:    "output",
		refName: "output",
		hasKey:  false,
	}
	BlockTypeModule = BlockType{
		name:    "module",
		refName: "module",
		hasKey:  true,
	}
	BlockTypeTerraform = BlockType{
		name:    "terraform",
		refName: "terraform",
		hasKey:  true,
	}
	BlockTypeUnknown = BlockType{
		name: "Unknown",
	}
)

var blockTypes = []BlockType{
	BlockTypeResource,
	BlockTypeVariable,
	BlockTypeData,
	BlockTypeLocal,
	BlockTypeProvider,
	BlockTypeOutput,
	BlockTypeModule,
	BlockTypeTerraform,
}

func GetBlockTypeByType(blockTypeStr string) BlockType {
	for _, bt := range blockTypes {
		if bt.name == blockTypeStr {
			return bt
		}
	}
	return BlockTypeUnknown
}

func getVariableValue(block Block) (*Attribute, error) {
	for _, attr := range block.Attributes {
		if attr.Name == "default" {
			return &attr, nil
		}
	}
	return nil, fmt.Errorf("could not find attribute")
}
