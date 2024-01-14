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

func newReference(parts []string) (*Reference, error) {
	var ref Reference

	if len(parts) == 0 {
		return nil, fmt.Errorf("cannot create empty reference")
	}
	bType := getBlockTypeByRef(parts[0])
	if bType.name != "Unknown" {
		ref.blockType = bType
		if len(parts) > 1 && ref.blockType.name != "resource" {
			ref.labels = parts[1:]
		}
	} else {
		ref.blockType = BlockTypeResource
		ref.labels = parts
	}

	if ref.blockType.hasKey {
		ref.key = ref.labels[len(ref.labels)-1]
		ref.labels = ref.labels[:len(ref.labels)-1]
	}

	ref.labels = append([]string{ref.blockType.name}, ref.labels...)

	return &ref, nil
}

type Type struct {
	name                  string
	refName               string
	removeTypeInReference bool
}

func (t Type) Name() string {
	return t.name
}

func (t Type) ShortName() string {
	if t.refName != "" {
		return t.refName
	}
	return t.name
}

func getBlockTypeByRef(blockTypeStr string) BlockType {
	for _, bt := range blockTypes {
		if bt.refName == blockTypeStr {
			return bt
		}
	}
	return BlockTypeUnknown
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
