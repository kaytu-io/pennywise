package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
	"strings"
)

type Block struct {
	Type        string
	Labels      []string
	Body        hcl.Body
	ChildBlocks []Block
	Attributes  []Attribute
	Context     *hcl.EvalContext
	CtxVariable cty.Value

	logger *zap.Logger
}

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
var (
	missingAttributeDiagnostic        = "Unsupported attribute"
	valueDoesNotHaveAnyIndices        = "Invalid index"
	valueIsNonIterableDiagnostic      = "Iteration over non-iterable value"
	invalidFunctionArgumentDiagnostic = "Invalid function argument"
)

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
				childBlocks, err := makeBlocks(logger, context, nil, &body.Blocks)
				if err != nil {
					return nil, err
				}
				attributes := make(hcl.Attributes)
				for _, a := range body.Attributes {
					attributes[a.Name] = a.AsHCLAttribute()
				}
				newBlock := Block{
					Type:        b.Type,
					Labels:      b.Labels,
					Body:        b.Body,
					ChildBlocks: childBlocks,
					Context:     context,
					logger:      logger,
				}
				newBlock.buildAttributes(attributes)
				if err != nil {
					return nil, err
				}
				totalBlocks = append(totalBlocks, newBlock)
			}
		}
	} else if childBlocks != nil {
		for _, b := range *childBlocks {
			childBlocks, err := makeBlocks(logger, context, nil, &b.Body.Blocks)
			if err != nil {
				return nil, err
			}
			attributes, diags := b.Body.JustAttributes()
			if diags.HasErrors() {
				return nil, diags
			}
			newBlock := Block{
				Type:        b.Type,
				Labels:      b.Labels,
				Body:        b.Body,
				ChildBlocks: childBlocks,
				Context:     context,
				logger:      logger,
			}
			newBlock.buildAttributes(attributes)
			if err != nil {
				return nil, err
			}
			totalBlocks = append(totalBlocks, newBlock)
		}
	}

	return totalBlocks, nil
}

func (b *Block) makeMapStructure(mappedBlocks map[string]interface{}) (map[string]interface{}, error) {
	mapStructure := make(map[string]interface{})
	ctxMapStructure := make(map[string]cty.Value)
	for _, attr := range b.Attributes {
		if attr.Name == "for_each" {
			continue
		}
		val, err := attr.Value()
		if attr.CtxVariable != nil {
			ctxMapStructure[attr.Name] = *attr.CtxVariable
		}
		if err != nil {
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
			blockValues, err := attrBlock.makeMapStructure(mappedBlocks)
			if err != nil {
				//b.logger.Error(fmt.Sprintf("error while getting %s value in block %s : %s", attr.Name, blockName, err.Error()))
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
			//b.logger.Debug(fmt.Sprintf("could not find attribute %s type in block %s", attr.Name, blockName))
		}
	}
	for _, childBlock := range b.ChildBlocks {
		var childBlockName string
		if len(childBlock.Labels) > 0 {
			childBlockName = fmt.Sprintf("%s.%s", childBlock.Type, strings.Join(childBlock.Labels, "."))
		} else {
			childBlockName = childBlock.Type
		}
		mappedChildBlock, err := childBlock.makeMapStructure(mappedBlocks)
		if err != nil {
			//b.logger.Error(fmt.Sprintf("error while making %s child block map structure in block %s : %s", childBlockName, blockName, err.Error()))
			continue
		}
		ctxMapStructure[childBlock.Type] = childBlock.CtxVariable
		mapStructure[childBlockName] = mappedChildBlock
	}
	ctxMapStructure["id"] = cty.EmptyObjectVal
	b.CtxVariable = cty.ObjectVal(ctxMapStructure)
	err := b.checkForEach()
	if err != nil {
		//b.logger.Error(fmt.Sprintf("error while parsing for each on block %s : %s", blockName, err.Error()))
	}
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

func (b *Block) checkForEach() error {
	forEach := b.findAttribute("for_each")
	if forEach == nil {
		return nil
	}

	ctyVal, diag := forEach.HclAttribute.Expr.Value(b.Context)
	if diag.HasErrors() {
		return nil
	}
	fmt.Println("foreach ctyVal", ctyVal.AsValueMap())
	return nil
}
