package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
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

	logger *zap.Logger
}

func makeNewBlock(logger *zap.Logger, ctx *hcl.EvalContext, b any, childHclBlocks *hclsyntax.Blocks, attributes hcl.Attributes) (*Block, error) {
	childBlocks, err := makeBlocks(logger, ctx, nil, childHclBlocks)
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
			logger:      logger,
		}
	} else if block, ok := b.(*hclsyntax.Block); ok {
		newBlock = Block{
			Type:        block.Type,
			Labels:      block.Labels,
			Body:        block.Body,
			ChildBlocks: childBlocks,
			Context:     ctx,
			logger:      logger,
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
	newBlock.buildAttributes(attributes)
	return &newBlock, nil
}

func (b *Block) cloneBlock(key string) Block {
	newBlock := *b
	newBlock.Name = fmt.Sprintf("%s[%s]", b.Name, key)
	return newBlock
}

func getFileBlocks(logger *zap.Logger, context *hcl.EvalContext, file *hcl.File) ([]Block, error) {
	contents, _, diags := file.Body.PartialContent(terraformSchema)
	if diags.HasErrors() {
		return nil, diags
	}
	myBlocks, err := makeBlocks(logger, context, &contents.Blocks, nil)
	if err != nil {
		return nil, err
	}
	return myBlocks, nil
}

func makeBlocks(logger *zap.Logger, context *hcl.EvalContext, blocks *hcl.Blocks, childBlocks *hclsyntax.Blocks) ([]Block, error) {
	var totalBlocks []Block
	if blocks != nil {
		for _, b := range *blocks {
			if body, ok := b.Body.(*hclsyntax.Body); ok {
				attributes := make(hcl.Attributes)
				for _, a := range body.Attributes {
					attributes[a.Name] = a.AsHCLAttribute()
				}
				newBlock, err := makeNewBlock(logger, context, b, &body.Blocks, attributes)
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
			newBlock, err := makeNewBlock(logger, context, b, &b.Body.Blocks, attributes)
			if err != nil {
				return nil, err
			}
			totalBlocks = append(totalBlocks, *newBlock)
		}
	}

	return totalBlocks, nil
}

func (b *Block) makeMapStructure(blockName string, ctx *hcl.EvalContext) (map[string]interface{}, error) {
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
			b.Diags.ChildDiags = append(b.Diags.ChildDiags, &attr.Diags)
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
			blockValues, err := attrBlock.makeMapStructure(attrBlockName, ctx)
			if err != nil {
				attr.Diags.Errors = append(attr.Diags.Errors, err)
				continue
			}
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
		mappedChildBlock, err := childBlock.makeMapStructure(childBlockName, ctx)
		if err != nil {
			continue
		}
		ctxMapStructure[childBlock.Type] = childBlock.CtxVariable
		mapStructure[childBlockName] = mappedChildBlock
	}
	ctxMapStructure["id"] = cty.StringVal(fmt.Sprintf("!ref:%s", blockName))
	b.CtxVariable = cty.ObjectVal(ctxMapStructure)

	return mapStructure, nil
}

func (b *Block) findAttribute(name string) *Attribute {
	for _, attr := range b.Attributes {
		if attr.Name == name {
			return &attr
		}
	}
	return nil
}

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
