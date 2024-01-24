package hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"strings"
)

// Attribute represents block attributes
type Attribute struct {
	Name         string
	HclAttribute hcl.Attribute
	Context      *hcl.EvalContext
	CtxVariable  *cty.Value
	Diags        Diags
}

// createAttributes creates list of attribute object by hcl.Attributes
func (b *Block) createAttributes(hclAttributes hcl.Attributes) {
	var attributes []Attribute
	for _, attr := range hclAttributes {
		attributes = append(attributes, Attribute{
			Name:         attr.Name,
			HclAttribute: *attr,
			Context:      b.Context,
			Diags:        Diags{Name: attr.Name, Type: AttributeDiag},
		})
	}
	b.Attributes = attributes
	return
}

// Value returns Attribute value by getting hcl.EvalContext
// uses the propagated context if ctx is nil
func (attr *Attribute) Value(ctx *hcl.EvalContext) (any, error) {
	if ctx == nil {
		ctx = attr.Context
	}
	ctyVal, diag := attr.HclAttribute.Expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	attr.CtxVariable = &ctyVal
	return getCtyValue(ctyVal)
}

func getCtyValue(ctyVal cty.Value) (any, error) {
	if isList(ctyVal) {
		return getListValues(ctyVal)
	}
	switch t := ctyVal.Type(); t {
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
	case cty.DynamicPseudoType:
		return nil, nil
	default:
		if checkMapCtyValue(ctyVal) {
			valueMap := make(map[string]any)
			for k, v := range ctyVal.AsValueMap() {
				value, err := getCtyValue(v)
				if err != nil {
					return nil, fmt.Errorf("unknown attribute type while getting value")
				}
				valueMap[k] = value
			}
			return valueMap, nil
		}
		return nil, fmt.Errorf("unknown attribute type while getting value")
	}
}

// isList checks if an attribute value is a kind of list or not
func isList(v cty.Value) bool {
	sourceTy := v.Type()

	return sourceTy.IsTupleType() || sourceTy.IsListType() || sourceTy.IsSetType()
}

// getListValues returns an Attribute value for the list values
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
			return nil, nil
		}
	} else {
		return nil, nil
	}
}

func checkMapCtyValue(value cty.Value) (b bool) {
	b = true
	defer func() {
		if r := recover(); r != nil {
			b = false
		}
	}()

	if !value.CanIterateElements() {
		b = false
	}
	_ = value.AsValueMap()

	return
}

func parseCtyValue(ctx *hcl.EvalContext, mapStructure map[string]interface{}, name string, val any) (map[string]interface{}, error) {
	switch val.(type) {
	case int64, int32, int:
		mapStructure[name] = val.(int64)
	case string:
		mapStructure[name] = val.(string)
	case bool:
		mapStructure[name] = val.(bool)
	case Block:
		attrBlock := val.(Block)
		var attrBlockName string
		if len(attrBlock.Labels) > 0 {
			attrBlockName = fmt.Sprintf("%s.%s", attrBlock.Type, strings.Join(attrBlock.Labels, "."))
		} else {
			attrBlockName = attrBlock.Type
		}
		blockValues := attrBlock.makeMapStructure(attrBlockName, ctx)
		mapStructure[name] = blockValues
	case []string:
		mapStructure[name] = val.([]string)
	case []bool:
		mapStructure[name] = val.([]bool)
	case []int64, []int32, []int:
		mapStructure[name] = val.([]int64)
	case map[string]any:
		valueMap := make(map[string]any)
		var err error
		for k, v := range val.(map[string]any) {
			valueMap, err = parseCtyValue(ctx, valueMap, k, v)
			if err != nil {
				return mapStructure, err
			}
		}
		mapStructure[name] = valueMap
	default:
		return mapStructure, fmt.Errorf("unknown attribute type while parsing")
	}
	return mapStructure, nil
}
