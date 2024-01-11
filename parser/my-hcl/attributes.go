package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Attribute struct {
	Name         string
	HclAttribute hcl.Attribute
}

func BuildAttributes(hclAttributes hcl.Attributes) []Attribute {
	var attributes []Attribute
	for _, attr := range hclAttributes {
		attributes = append(attributes, Attribute{
			Name:         attr.Name,
			HclAttribute: *attr,
		})
	}
	return attributes
}

func (attr *Attribute) readAttributeReference() (*Reference, error) {
	expression := attr.HclAttribute.Expr
	switch t := expression.(type) {
	case *hclsyntax.FunctionCallExpr:
		fmt.Println(t)
		fmt.Println("FunctionCallExpr")
	case *hclsyntax.ConditionalExpr:
		fmt.Println("ConditionalExpr")
	case *hclsyntax.ScopeTraversalExpr:
		var refParts []string

		for _, x := range t.Variables() {
			for _, p := range x {
				switch part := p.(type) {
				case hcl.TraverseRoot:
					refParts = append(refParts, part.Name)
				case hcl.TraverseAttr:
					refParts = append(refParts, part.Name)
				case hcl.TraverseIndex:
					refParts[len(refParts)-1] = fmt.Sprintf("%s[%s]", refParts[len(refParts)-1], attr.getIndexValue(part))
				}
			}
		}
		return newReference(refParts)
	case *hclsyntax.TemplateWrapExpr:
		fmt.Println("TemplateWrapExpr")
	case *hclsyntax.TemplateExpr:
		fmt.Println("TemplateExpr")
	case *hclsyntax.TupleConsExpr:
		fmt.Println("TupleConsExpr")
	case *hclsyntax.RelativeTraversalExpr:
		fmt.Println("RelativeTraversalExpr")
	case *hclsyntax.IndexExpr:
		fmt.Println("IndexExpr")
	default:
		fmt.Println("DEFAULT")
	}
	return nil, fmt.Errorf("unknown type")
}

func (attr *Attribute) getIndexValue(part hcl.TraverseIndex) string {
	switch part.Key.Type() {
	case cty.String:
		return fmt.Sprintf("%q", part.Key.AsString())
	case cty.Number:
		var intVal int
		if err := gocty.FromCtyValue(part.Key, &intVal); err != nil {
			return "0"
		}
		return fmt.Sprintf("%d", intVal)
	default:
		return "0"
	}
}

func (attr *Attribute) Value(mappedBlocks map[string]interface{}) (any, error) {
	ctx := hcl.EvalContext{}
	ctx.Variables = make(map[string]cty.Value)
	ctyVal, diag := attr.HclAttribute.Expr.Value(&ctx)
	if diag.HasErrors() {
		if diag.HasErrors() {
			ref, err := attr.readAttributeReference()
			if err != nil {
				return nil, err
			}
			return getRefValue(mappedBlocks, *ref)
		}
	}
	if isList(ctyVal) {
		return getListValues(ctyVal)
	}
	switch ctyVal.Type() {
	case cty.String:
		var s string
		err := gocty.FromCtyValue(ctyVal, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	case cty.Number:
		var i int64
		err := gocty.FromCtyValue(ctyVal, &i)
		if err != nil {
			return nil, err
		}
		return i, nil
	case cty.Bool:
		var b bool
		err := gocty.FromCtyValue(ctyVal, &b)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, fmt.Errorf("value type not implemented")
	}
}

func isList(v cty.Value) bool {
	sourceTy := v.Type()

	return sourceTy.IsTupleType() || sourceTy.IsListType() || sourceTy.IsSetType()
}

func getListValues(ctyVal cty.Value) (any, error) {
	it := ctyVal.ElementIterator()
	if it.Next() {
		key, sourceItem := it.Element()
		switch sourceItem.Type() {
		case cty.String:
			items := make([]string, ctyVal.LengthInt())
			var s string
			err := gocty.FromCtyValue(sourceItem, &s)
			if err != nil {
				return nil, err
			}
			i, _ := key.AsBigFloat().Int64()
			items[i] = s
			for it.Next() {
				key, sourceItem = it.Element()
				var s string
				err = gocty.FromCtyValue(sourceItem, &s)
				if err != nil {
					return nil, err
				}
				i, _ = key.AsBigFloat().Int64()
				items[i] = s
			}
			return items, nil
		case cty.Number:
			items := make([]int64, ctyVal.LengthInt())
			var v int64
			err := gocty.FromCtyValue(sourceItem, &v)
			if err != nil {
				return nil, err
			}
			i, _ := key.AsBigFloat().Int64()
			items[i] = v
			for it.Next() {
				key, sourceItem = it.Element()
				var v int64
				err = gocty.FromCtyValue(sourceItem, &v)
				if err != nil {
					return nil, err
				}
				i, _ = key.AsBigFloat().Int64()
				items[i] = v
			}
			return items, nil
		case cty.Bool:
			items := make([]bool, ctyVal.LengthInt())
			var b bool
			err := gocty.FromCtyValue(sourceItem, &b)
			if err != nil {
				return nil, err
			}
			i, _ := key.AsBigFloat().Int64()
			items[i] = b
			for it.Next() {
				key, sourceItem = it.Element()
				var b bool
				err = gocty.FromCtyValue(sourceItem, &b)
				if err != nil {
					return nil, err
				}
				i, _ = key.AsBigFloat().Int64()
				items[i] = b
			}
			return items, nil
		default:
			return nil, fmt.Errorf("list value type not implemented")
		}
	} else {
		return nil, nil
	}
}

func getRefValue(mappedBlocks map[string]interface{}, reference Reference) (any, error) {
	if len(reference.labels) > 0 {
		block, err := findRefBlockFromLabels(mappedBlocks, reference.labels)
		if err != nil {
			return nil, err
		}
		if reference.blockType.hasKey {
			if reference.key == "id" || reference.key == "name" {
				return *block, nil
			}
			attr := findAttribute(*block, reference.key)
			if attr == nil {
				return nil, fmt.Errorf("could not find attribute")
			} else {
				value, err := attr.Value(mappedBlocks)
				if err != nil {
					return nil, err
				}
				return value, nil
			}
		} else {
			if reference.blockType.getValueFunction != nil {
				attr, err := reference.blockType.getValueFunction(*block)
				if err != nil {
					return nil, err
				}
				value, err := attr.Value(mappedBlocks)
				if err != nil {
					return nil, err
				}
				return value, nil
			}
			return nil, fmt.Errorf("not handled yet")
		}
	} else {
		_ = mappedBlocks[reference.blockType.name]
		return nil, fmt.Errorf("not handled yet")
	}
}

func findRefBlockFromLabels(mappedBlocks map[string]interface{}, labels []string) (*Block, error) {
	if len(labels) > 1 {
		labeledMappedBlocks := mappedBlocks[labels[0]]
		if _, ok := labeledMappedBlocks.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("wrong ref labels: %s", labels)
		}
		return findRefBlockFromLabels(labeledMappedBlocks.(map[string]interface{}), labels[1:])
	} else {
		block := mappedBlocks[labels[0]]
		if _, ok := block.(Block); !ok {
			return nil, fmt.Errorf("wrong ref labels: %s", labels)
		}
		result := block.(Block)
		return &result, nil
	}
}

func findAttribute(block Block, attrName string) *Attribute { // handle attributes in blocks
	for _, attr := range block.Attributes {
		if attr.Name == attrName {
			return &attr
		}
	}
	return nil
}
