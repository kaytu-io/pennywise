package my_hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// findBadVariablesFromExpression attempts to find the variables that are missing by calling the underlying expressions
// and checking if they have any missing attributes diagnostics. findBadVariablesFromExpression is a fallback method
// as normally Diagnostics return the variables we need. However, in cases where the expressions are complex (e.g.
// a splat expression within a function call) the Diagnostics will only have variable information from the last expression.
// Meaning that in many cases they won't actually contain the problem variables and calling diag.Variables() will return nil.
func (attr *Attribute) findBadVariablesFromExpression(expression hcl.Expression) []hcl.Traversal {
	var badVars []hcl.Traversal
	ctx := attr.Ctx.Inner()
	switch t := expression.(type) {
	case *hclsyntax.ForExpr:
		// if there are bad vars in the collection we need to evaluate these first
		badVars = append(badVars, attr.findBadVariablesFromExpression(t.CollExpr)...)
		if badVars != nil {
			return badVars
		}

		collVal, _ := t.CollExpr.Value(ctx)
		if !isList(collVal) {
			collVal = cty.TupleVal([]cty.Value{collVal})
		}
		it := collVal.ElementIterator()

		for it.Next() {
			k, v := it.Element()
			childCtx := ctx.NewChild()
			childCtx.Variables = map[string]cty.Value{}
			if t.KeyVar != "" {
				childCtx.Variables[t.KeyVar] = k
			}
			childCtx.Variables[t.ValVar] = v

			if t.CondExpr != nil {
				_, diags := t.CondExpr.Value(childCtx)
				if isAttrMissing(diags) {
					trav := findBadTraversal(childCtx, t.CondExpr)

					if trav.RootName() == t.ValVar {
						abs, _ := hcl.AbsTraversalForExpr(t.CollExpr)
						rels := toRelativeTraversal(trav)
						traversal := hcl.TraversalJoin(abs, rels)
						badVars = append(badVars, traversal)

						return badVars
					}
				}

			}
		}

		if t.CondExpr != nil {
			badVars = append(badVars, attr.findBadVariablesFromExpression(t.CondExpr)...)
		}

		badVars = append(badVars, attr.findBadVariablesFromExpression(t.ValExpr)...)
		return badVars
	case *hclsyntax.FunctionCallExpr:
		for _, arg := range t.Args {
			badVars = append(badVars, attr.findBadVariablesFromExpression(arg)...)
		}

		return badVars
	case *hclsyntax.ConditionalExpr:
		badVars = append(badVars, attr.findBadVariablesFromExpression(t.TrueResult)...)
		badVars = append(badVars, attr.findBadVariablesFromExpression(t.FalseResult)...)
		badVars = append(badVars, attr.findBadVariablesFromExpression(t.Condition)...)
		return badVars
	case *hclsyntax.TemplateWrapExpr:
		return attr.findBadVariablesFromExpression(t.Wrapped)
	case *hclsyntax.TemplateExpr:
		for _, part := range t.Parts {
			badVars = append(badVars, attr.findBadVariablesFromExpression(part)...)
		}

		return badVars
	case *hclsyntax.RelativeTraversalExpr:
		switch s := t.Source.(type) {
		case *hclsyntax.IndexExpr:
			ctx := attr.Ctx.Inner()
			val, diags := s.Collection.Value(ctx)
			if isAttrMissing(diags) {
				return attr.findBadVariables(s.Collection.Variables())
			}

			it := val.ElementIterator()
			for it.Next() {
				_, v := it.Element()

				_, d := t.Traversal.TraverseRel(v)
				if isAttrMissing(d) {
					traversal := s.Collection.Variables()[0]
					traversal = append(traversal, hcl.TraverseSplat{})
					traversal = append(traversal, t.Traversal...)
					badVars = append(badVars, traversal)
					return badVars
				} else {
					break
				}
			}

			return attr.findBadVariables(s.Collection.Variables())
		default:
			return attr.findBadVariablesFromExpression(t.Source)
		}
	case *hclsyntax.IndexExpr:
		return attr.findBadVariables(t.Collection.Variables())
	case *hclsyntax.SplatExpr:
		_, diag := t.Value(attr.Ctx.Inner())
		if isAttrMissing(diag) {
			baseVars := t.Variables()

			if rt, ok := t.Each.(*hclsyntax.RelativeTraversalExpr); ok {
				for i, baseVar := range baseVars {
					baseVars[i] = append(baseVar, rt.Traversal...)
				}

				badVars = append(badVars, baseVars...)
				return badVars
			}
		}
	case *hclsyntax.ObjectConsExpr:
		for _, item := range t.Items {
			badVars = append(badVars, attr.findBadVariablesFromExpression(item.KeyExpr)...)
			badVars = append(badVars, attr.findBadVariablesFromExpression(item.ValueExpr)...)
		}

		return badVars
	}

	return attr.findBadVariables(expression.Variables())
}

func isAttrMissing(diag hcl.Diagnostics) bool {
	for _, d := range diag {
		if d.Summary == missingAttributeDiagnostic {
			return true
		}
	}

	return false
}

func findBadTraversal(ctx *hcl.EvalContext, expression hcl.Expression) hcl.Traversal {
	switch t := expression.(type) {
	case *hclsyntax.ConditionalExpr:
		if b := findBadTraversal(ctx, t.TrueResult); b != nil {
			return b
		}
		if b := findBadTraversal(ctx, t.FalseResult); b != nil {
			return b
		}
		if b := findBadTraversal(ctx, t.Condition); b != nil {
			return b
		}
	case *hclsyntax.BinaryOpExpr:
		if b := findBadTraversal(ctx, t.LHS); b != nil {
			return b
		}
		if b := findBadTraversal(ctx, t.RHS); b != nil {
			return b
		}
	}

	_, diag := expression.Value(ctx)
	if isAttrMissing(diag) {
		trav, _ := hcl.AbsTraversalForExpr(expression)
		return trav
	}

	return nil
}

func (attr *Attribute) findBadVariables(traversals []hcl.Traversal) []hcl.Traversal {
	var badVars []hcl.Traversal

	for _, traversal := range traversals {
		if traversal.IsRelative() {
			continue
		}

		_, diag := traversal.TraverseAbs(attr.Ctx.Inner())
		if isAttrMissing(diag) {
			badVars = append(badVars, traversal)
		}
	}

	return badVars
}

func toRelativeTraversal(traversal hcl.Traversal) hcl.Traversal {
	var ret hcl.Traversal
	for _, traverser := range traversal {
		if _, ok := traverser.(hcl.TraverseRoot); ok {
			continue
		}

		ret = append(ret, traverser)
	}

	return ret
}
