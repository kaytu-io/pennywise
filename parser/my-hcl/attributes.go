package my_hcl

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"go.uber.org/zap"
)

type Attribute struct {
	Name         string
	HclAttribute hcl.Attribute
	Context      *hcl.EvalContext
	CtxVariable  *cty.Value

	logger *zap.Logger
}

func (b *Block) buildAttributes(hclAttributes hcl.Attributes) {
	var attributes []Attribute
	for _, attr := range hclAttributes {
		attributes = append(attributes, Attribute{
			Name:         attr.Name,
			HclAttribute: *attr,
			Context:      b.Context,
			logger:       b.logger,
		})
	}
	b.Attributes = attributes
	return
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

func (attr *Attribute) Value() (any, error) {
	ctyVal, diag := attr.HclAttribute.Expr.Value(attr.Context)
	if diag.HasErrors() {
		fmt.Println("ERROR", attr.Name, diag[0].Detail)
		return nil, nil
	}
	attr.CtxVariable = &ctyVal
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
		fmt.Println("Unknown attribute", attr.Name, attr.HclAttribute)
		return nil, nil
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
			return nil, nil
		}
	} else {
		return nil, nil
	}
}
