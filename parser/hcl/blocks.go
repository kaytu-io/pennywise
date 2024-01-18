package hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"strings"
)

var (
	terraformSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "terraform",
			},
			{
				Type:       "provider",
				LabelNames: []string{"name"},
			},
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type: "locals",
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
			},
		},
	}
)

// BlockType represents a block type
type BlockType struct {
	name    string
	refName string
}

var (
	BlockTypeResource = BlockType{
		name:    "resource",
		refName: "",
	}
	BlockTypeVariable = BlockType{
		name:    "variable",
		refName: "var",
	}
	BlockTypeData = BlockType{
		name:    "data",
		refName: "data",
	}
	BlockTypeLocal = BlockType{
		name:    "locals",
		refName: "local",
	}
	BlockTypeProvider = BlockType{
		name:    "provider",
		refName: "provider",
	}
	BlockTypeOutput = BlockType{
		name:    "output",
		refName: "output",
	}
	BlockTypeModule = BlockType{
		name:    "module",
		refName: "module",
	}
	BlockTypeTerraform = BlockType{
		name:    "terraform",
		refName: "terraform",
	}
	BlockTypeUnknown = BlockType{
		name: "Unknown",
	}
)

// blockTypes available block types
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

// GetBlockTypeByType returns BlockType by name
func GetBlockTypeByType(blockTypeStr string) BlockType {
	for _, bt := range blockTypes {
		if bt.name == blockTypeStr {
			return bt
		}
	}
	return BlockTypeUnknown
}

// Block represents a block in a hcl file
type Block struct {
	Name        string
	Type        string
	Labels      []string
	Body        hcl.Body
	ChildBlocks []Block
	Attributes  []Attribute
	Context     *hcl.EvalContext
	CtxVariable cty.Value
	Diags       Diags
}

// makeNewBlock creates new Block by hcl syntax block
func makeNewBlock(ctx *hcl.EvalContext, b any, childHclBlocks *hclsyntax.Blocks, attributes hcl.Attributes) (*Block, error) {
	childBlocks, err := makeBlocks(ctx, nil, childHclBlocks)
	if err != nil {
		return nil, err
	}
	var newBlock Block
	if block, ok := b.(*hcl.Block); ok {
		newBlock = Block{
			Type:        block.Type,
			Labels:      block.Labels,
			Body:        block.Body,
			ChildBlocks: childBlocks,
			Context:     ctx,
		}
	} else if block, ok := b.(*hclsyntax.Block); ok {
		newBlock = Block{
			Type:        block.Type,
			Labels:      block.Labels,
			Body:        block.Body,
			ChildBlocks: childBlocks,
			Context:     ctx,
		}
	} else {
		return nil, fmt.Errorf("invalid block")
	}

	var blockName string
	if len(newBlock.Labels) > 0 {
		blockName = fmt.Sprintf("%s.%s", newBlock.Type, strings.Join(newBlock.Labels, "."))
	} else {
		blockName = newBlock.Type
	}
	newBlock.Name = blockName
	newBlock.Diags = Diags{Name: newBlock.Name, Type: BlockDiag}
	newBlock.createAttributes(attributes)
	return &newBlock, nil
}

// getFileBlocks returns list of Block in a hcl file
func getFileBlocks(context *hcl.EvalContext, file *hcl.File) ([]Block, error) {
	contents, _, diags := file.Body.PartialContent(terraformSchema)
	if diags.HasErrors() {
		return nil, diags
	}
	myBlocks, err := makeBlocks(context, &contents.Blocks, nil)
	if err != nil {
		return nil, err
	}
	return myBlocks, nil
}

// makeBlocks create list of Block by getting list of hcl blocks
func makeBlocks(context *hcl.EvalContext, blocks *hcl.Blocks, childBlocks *hclsyntax.Blocks) ([]Block, error) {
	var totalBlocks []Block
	if blocks != nil {
		for _, b := range *blocks {
			if body, ok := b.Body.(*hclsyntax.Body); ok {
				attributes := make(hcl.Attributes)
				for _, a := range body.Attributes {
					attributes[a.Name] = a.AsHCLAttribute()
				}
				newBlock, err := makeNewBlock(context, b, &body.Blocks, attributes)
				if err != nil {
					return nil, err
				}
				totalBlocks = append(totalBlocks, *newBlock)
			}
		}
	} else if childBlocks != nil {
		for _, b := range *childBlocks {
			attributes, diags := b.Body.JustAttributes()
			if diags.HasErrors() {
				return nil, diags
			}
			newBlock, err := makeNewBlock(context, b, &b.Body.Blocks, attributes)
			if err != nil {
				return nil, err
			}
			totalBlocks = append(totalBlocks, *newBlock)
		}
	}

	return totalBlocks, nil
}

// makeMapStructure makes a map structure for a block
func (b *Block) makeMapStructure(blockName string, ctx *hcl.EvalContext) map[string]interface{} {
	if ctx == nil {
		ctx = b.Context
	}
	mapStructure := make(map[string]interface{})
	ctxMapStructure := make(map[string]cty.Value)
	for _, attr := range b.Attributes {
		if attr.Name == "for_each" {
			continue
		}
		val, err := attr.Value(ctx)
		if attr.CtxVariable != nil {
			ctxMapStructure[attr.Name] = *attr.CtxVariable
		}
		if err != nil {
			attr.Diags.Errors = append(attr.Diags.Errors, err)
			b.Diags.ChildDiags = append(b.Diags.ChildDiags, attr.Diags)
			continue
		}
		switch val.(type) {
		case int64, int32, int:
			mapStructure[attr.Name] = val.(int64)
		case string:
			mapStructure[attr.Name] = val.(string)
		case bool:
			mapStructure[attr.Name] = val.(bool)
		case Block:
			attrBlock := val.(Block)
			var attrBlockName string
			if len(attrBlock.Labels) > 0 {
				attrBlockName = fmt.Sprintf("%s.%s", attrBlock.Type, strings.Join(attrBlock.Labels, "."))
			} else {
				attrBlockName = attrBlock.Type
			}
			blockValues := attrBlock.makeMapStructure(attrBlockName, ctx)
			mapStructure[attr.Name] = blockValues
		case []string:
			mapStructure[attr.Name] = val.([]string)
		case []bool:
			mapStructure[attr.Name] = val.([]bool)
		case []int64, []int32, []int:
			mapStructure[attr.Name] = val.([]int64)
		default:
			attr.Diags.Errors = append(attr.Diags.Errors, fmt.Errorf("unknown attribute type while parsing"))
		}
	}
	for _, childBlock := range b.ChildBlocks {
		var childBlockName string
		if len(childBlock.Labels) > 0 {
			childBlockName = fmt.Sprintf("%s.%s", childBlock.Type, strings.Join(childBlock.Labels, "."))
		} else {
			childBlockName = childBlock.Type
		}
		mappedChildBlock := childBlock.makeMapStructure(childBlockName, ctx)
		ctxMapStructure[childBlock.Type] = childBlock.CtxVariable
		if _, ok := mapStructure[childBlockName]; !ok {
			mapStructure[childBlockName] = make([]map[string]interface{}, 0)
		}
		mapStructure[childBlockName] = append(mapStructure[childBlockName].([]map[string]interface{}), mappedChildBlock)
	}
	ctxMapStructure["id"] = cty.StringVal(fmt.Sprintf("%s", blockName))
	mapStructure["id"] = fmt.Sprintf("%s", blockName)
	b.CtxVariable = cty.ObjectVal(ctxMapStructure)

	return mapStructure
}

// findAttribute find an Attribute in a block by name
func (b *Block) findAttribute(name string) *Attribute {
	for _, attr := range b.Attributes {
		if attr.Name == name {
			return &attr
		}
	}
	return nil
}

// checkForEach returns for_each values if the expression exists in the Block
func (b *Block) checkForEach() (map[string]cty.Value, error) {
	forEach := b.findAttribute("for_each")
	if forEach == nil {
		return nil, nil
	}

	ctyVal, diag := forEach.HclAttribute.Expr.Value(b.Context)
	if diag.HasErrors() {
		return nil, diag
	}
	return ctyVal.AsValueMap(), nil
}
