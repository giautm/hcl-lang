package decoder

import (
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type ExprConstraints schema.ExprConstraints

func (ec ExprConstraints) FriendlyNames() []string {
	names := make([]string, 0)
	for _, constraint := range ec {
		if name := constraint.FriendlyName(); name != "" &&
			!namesContain(names, name) {
			names = append(names, name)
		}
	}
	return names
}

func namesContain(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func (ec ExprConstraints) HasKeywordsOnly() bool {
	hasKeywordExpr := false
	for _, constraint := range ec {
		if _, ok := constraint.(schema.KeywordExpr); ok {
			hasKeywordExpr = true
		} else {
			return false
		}
	}
	return hasKeywordExpr
}

func (ec ExprConstraints) KeywordExpr() (schema.KeywordExpr, bool) {
	for _, c := range ec {
		if kw, ok := c.(schema.KeywordExpr); ok {
			return kw, ok
		}
	}
	return schema.KeywordExpr{}, false
}

func (ec ExprConstraints) MapExpr() (schema.MapExpr, bool) {
	for _, c := range ec {
		if me, ok := c.(schema.MapExpr); ok {
			return me, ok
		}
	}
	return schema.MapExpr{}, false
}

func (ec ExprConstraints) TupleConsExpr() (schema.TupleConsExpr, bool) {
	for _, c := range ec {
		if tc, ok := c.(schema.TupleConsExpr); ok {
			return tc, ok
		}
	}
	return schema.TupleConsExpr{}, false
}

func (ec ExprConstraints) HasLiteralTypeOf(exprType cty.Type) bool {
	for _, c := range ec {
		if lt, ok := c.(schema.LiteralTypeExpr); ok && lt.Type.Equals(exprType) {
			return true
		}
	}
	return false
}

func (ec ExprConstraints) HasLiteralValueOf(val cty.Value) bool {
	for _, c := range ec {
		if lv, ok := c.(schema.LiteralValue); ok && lv.Val.RawEquals(val) {
			return true
		}
	}
	return false
}

func (ec ExprConstraints) LiteralValueOf(val cty.Value) (schema.LiteralValue, bool) {
	for _, c := range ec {
		if lv, ok := c.(schema.LiteralValue); ok && lv.Val.RawEquals(val) {
			return lv, true
		}
	}
	return schema.LiteralValue{}, false
}

func (ec ExprConstraints) LiteralTypeOfTupleExpr() (schema.LiteralTypeExpr, bool) {
	for _, c := range ec {
		if lv, ok := c.(schema.LiteralTypeExpr); ok {
			if lv.Type.IsListType() {
				return lv, true
			}
			if lv.Type.IsSetType() {
				return lv, true
			}
			if lv.Type.IsTupleType() {
				return lv, true
			}
		}
	}
	return schema.LiteralTypeExpr{}, false
}

func (ec ExprConstraints) LiteralTypeOfObjectConsExpr() (schema.LiteralTypeExpr, bool) {
	for _, c := range ec {
		if lv, ok := c.(schema.LiteralTypeExpr); ok {
			if lv.Type.IsObjectType() {
				return lv, true
			}
			if lv.Type.IsMapType() {
				return lv, true
			}
		}
	}
	return schema.LiteralTypeExpr{}, false
}

func (ec ExprConstraints) LiteralValueOfTupleExpr(expr *hclsyntax.TupleConsExpr) (schema.LiteralValue, bool) {
	exprValues := make([]cty.Value, len(expr.Exprs))
	for i, e := range expr.Exprs {
		val, _ := e.Value(nil)
		if !val.IsWhollyKnown() || val.IsNull() {
			return schema.LiteralValue{}, false
		}
		exprValues[i] = val
	}

	for _, c := range ec {
		if lv, ok := c.(schema.LiteralValue); ok {
			valType := lv.Val.Type()
			if valType.IsListType() && lv.Val.RawEquals(cty.ListVal(exprValues)) {
				return lv, true
			}
			if valType.IsSetType() && lv.Val.RawEquals(cty.SetVal(exprValues)) {
				return lv, true
			}
			if valType.IsTupleType() && lv.Val.RawEquals(cty.TupleVal(exprValues)) {
				return lv, true
			}
		}
	}

	return schema.LiteralValue{}, false
}

func (ec ExprConstraints) LiteralValueOfObjectConsExpr(expr *hclsyntax.ObjectConsExpr) (schema.LiteralValue, bool) {
	exprValues := make(map[string]cty.Value)
	for _, item := range expr.Items {
		key, _ := item.KeyExpr.Value(nil)
		if !key.IsWhollyKnown() || key.Type() != cty.String {
			return schema.LiteralValue{}, false
		}

		val, _ := item.ValueExpr.Value(nil)
		if !val.IsWhollyKnown() || val.IsNull() {
			return schema.LiteralValue{}, false
		}

		exprValues[key.AsString()] = val
	}

	for _, c := range ec {
		if lv, ok := c.(schema.LiteralValue); ok {
			valType := lv.Val.Type()
			if valType.IsMapType() && lv.Val.RawEquals(cty.MapVal(exprValues)) {
				return lv, true
			}
			if valType.IsObjectType() && lv.Val.RawEquals(cty.ObjectVal(exprValues)) {
				return lv, true
			}
		}
	}

	return schema.LiteralValue{}, false
}
