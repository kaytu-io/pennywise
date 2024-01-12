package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"strings"
)

type Block struct {
	Type        string
	Labels      []string
	Body        hcl.Body
	ChildBlocks []Block
	Attributes  []Attribute
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

func getFileBlocks(file *hcl.File) ([]Block, error) {
	contents, _, diags := file.Body.PartialContent(terraformSchema)
	if diags.HasErrors() {
		return nil, diags
	}
	myBlocks, err := makeBlocks(&contents.Blocks, nil)
	if err != nil {
		return nil, err
	}
	return myBlocks, nil
}

func makeBlocks(blocks *hcl.Blocks, childBlocks *hclsyntax.Blocks) ([]Block, error) {
	var totalBlocks []Block
	if blocks != nil {
		for _, b := range *blocks {
			if body, ok := b.Body.(*hclsyntax.Body); ok {
				childBlocks, err := makeBlocks(nil, &body.Blocks)
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
					Attributes:  BuildAttributes(attributes),
				}
				err = newBlock.checkForEach()
				if err != nil {
					return nil, err
				}
				totalBlocks = append(totalBlocks, newBlock)
			}
		}
	} else if childBlocks != nil {
		for _, b := range *childBlocks {
			childBlocks, err := makeBlocks(nil, &b.Body.Blocks)
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
				Attributes:  BuildAttributes(attributes),
			}
			err = newBlock.checkForEach()
			if err != nil {
				return nil, err
			}
			totalBlocks = append(totalBlocks, newBlock)
		}
	}

	return totalBlocks, nil
}

func (b Block) MakeMapStructure(mappedBlocks map[string]interface{}) (map[string]interface{}, error) {
	mapStructure := make(map[string]interface{})
	for _, attr := range b.Attributes {
		val, err := attr.Value(mappedBlocks)
		if err != nil {
			return nil, err
		}
		switch val.(type) {
		case int64, int32, int:
			mapStructure[attr.Name] = val.(int64)
		case string:
			mapStructure[attr.Name] = val.(string)
		case bool:
			mapStructure[attr.Name] = val.(bool)
		case Block:
			blockValues, err := val.(Block).MakeMapStructure(mappedBlocks)
			if err != nil {
				return nil, err
			}
			mapStructure[attr.Name] = blockValues
		case []string:
			mapStructure[attr.Name] = val.([]string)
		case []bool:
			mapStructure[attr.Name] = val.([]bool)
		case []int64, []int32, []int:
			mapStructure[attr.Name] = val.([]int64)
		default:
			return nil, nil
		}
	}
	for _, childBlock := range b.ChildBlocks {
		var blockName string
		if len(childBlock.Labels) > 0 {
			blockName = fmt.Sprintf("%s.%s", childBlock.Type, strings.Join(childBlock.Labels, "."))
		} else {
			blockName = childBlock.Type
		}
		mappedChildBlock, err := childBlock.MakeMapStructure(mappedBlocks)
		if err != nil {
			return nil, err
		}
		mapStructure[blockName] = mappedChildBlock
	}
	return mapStructure, nil
}

func (b Block) findAttribute(name string) *Attribute {
	for _, attr := range b.Attributes {
		if attr.Name == name {
			return &attr
		}
	}
	return nil
}

func (b Block) checkForEach() error {
	forEach := b.findAttribute("for_each")
	if forEach == nil {
		return nil
	}
	ctx := hcl.EvalContext{}
	ctx.Variables = make(map[string]cty.Value)
	//ctx.Variables["local"] = cty.ObjectVal(map[string]cty.Value{"permutations": cty.ObjectVal(make(map[string]cty.Value))})
	//fmt.Println(ctx.Variables)

	for _, traversal := range forEach.HclAttribute.Expr.Variables() {
		var rootName string
		for _, traverser := range traversal {
			fmt.Println("Traversal", traversal)
			if r, ok := traverser.(hcl.TraverseRoot); ok {
				rootName = r.Name
				break
			}
		}

		ob := ctx.Variables[rootName]
		fmt.Println("pb", ob)
		if ob.IsNull() || !ob.IsKnown() {
			ob = cty.ObjectVal(make(map[string]cty.Value))
		}

		ctx.Variables[rootName] = buildObject(traversal, ob, 0)
		fmt.Println(ctx.Variables)
	}
	ctyVal, diag := forEach.HclAttribute.Expr.Value(&ctx)
	if diag.HasErrors() {
		for _, d := range diag {
			fmt.Println("diag", d.Detail)
			if d.Summary == missingAttributeDiagnostic || d.Summary == valueDoesNotHaveAnyIndices {
				fmt.Println(missingAttributeDiagnostic, valueDoesNotHaveAnyIndices)
			}
			if d.Summary == valueIsNonIterableDiagnostic {
				fmt.Println(valueIsNonIterableDiagnostic)
			}
			if d.Summary == invalidFunctionArgumentDiagnostic {
				fmt.Println(invalidFunctionArgumentDiagnostic)
			}
		}
		return diag
	}
	fmt.Println("ctyVal", ctyVal)
	ctyVal.ForEachElement(func(key cty.Value, val cty.Value) bool {
		fmt.Println("here2")
		fmt.Println("key:", key)
		fmt.Println("val:", val)
		return false
	})
	return nil
}

func buildObject(traversal hcl.Traversal, value cty.Value, i int) cty.Value {
	if i > len(traversal)-1 {
		return value
	}
	traverser := traversal[i]

	var valueMap map[string]cty.Value

	if value.IsKnown() && !value.IsNull() && value.CanIterateElements() {
		valueMap = value.AsValueMap()
	}

	if valueMap == nil {
		valueMap = make(map[string]cty.Value)
	}

	// traverse splat is a special holding type which means we want to traverse all the attributes on the map.
	if _, ok := traverser.(hcl.TraverseSplat); ok {
		for k, v := range valueMap {
			if v.Type().IsObjectType() {
				valueMap[k] = buildObject(traversal, v, i+1)
				continue
			}

			valueMap[k] = v
		}

		return cty.ObjectVal(valueMap)
	}

	if index, ok := traverser.(hcl.TraverseIndex); ok {
		kc, err := convert.Convert(index.Key, cty.String)
		if err != nil {
			kc = cty.StringVal("0")
		}

		k := kc.AsString()

		if vv, exists := valueMap[k]; exists {
			valueMap[k] = buildObject(traversal, vv, i+1)
			return cty.ObjectVal(valueMap)
		}

		valueMap[k] = buildObject(traversal, cty.ObjectVal(make(map[string]cty.Value)), i+1)

		return cty.ObjectVal(valueMap)
	}

	if v, ok := traverser.(hcl.TraverseAttr); ok {
		fmt.Println("TRAVERSER NAME", v.Name)
		if len(traversal)-1 == i {
			// if the attribute already exists, and we're not setting a list value
			// then we should return here. It's most likely that we weren't able to
			// get the full variable calls for the context, so resetting the value could
			// be harmful.
			if _, exists := valueMap[v.Name]; exists {
				return value
			}
			valueMap[v.Name] = cty.ObjectVal(make(map[string]cty.Value))
			return cty.ObjectVal(valueMap)
		}
		if vv, exists := valueMap[v.Name]; exists {
			if isList(vv) {
				items := make([]cty.Value, vv.LengthInt())
				it := vv.ElementIterator()
				for it.Next() {
					key, sourceItem := it.Element()
					val := buildObject(traversal, sourceItem, i+1)
					i, _ := key.AsBigFloat().Int64()
					items[i] = val
				}
				valueMap[v.Name] = cty.TupleVal(items)
				return cty.ObjectVal(valueMap)
			}

			valueMap[v.Name] = buildObject(traversal, vv, i+1)
			return cty.ObjectVal(valueMap)
		}

		valueMap[v.Name] = buildObject(traversal, cty.ObjectVal(make(map[string]cty.Value)), i+1)
		return cty.ObjectVal(valueMap)
	}
	return buildObject(traversal, value, i+1)
}
