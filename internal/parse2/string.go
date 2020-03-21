package parse2

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type String struct {
	Value     string
	attribute *hcl.Attribute
	val       cty.Value
	ctx       *hcl.EvalContext
}

func newStringAttribute(attr *hcl.Attribute, ctx *hcl.EvalContext) (sa *String, diag hcl.Diagnostics) {
	sa = &String{
		attribute: attr,
	}
	_, diag = sa.fill(ctx)

	return
}

func (a *String) fill(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	val, diag := a.attribute.Expr.Value(ctx)
	if diag.HasErrors() {
		return a.val, diag
	}

	t := val.Type()
	if t != cty.String {
		diag := hcl.Diagnostic{
			Summary:     "invalid type",
			Detail:      fmt.Sprintf("expected string, got %v", t.FriendlyName()),
			Severity:    hcl.DiagError,
			Subject:     &a.attribute.Range,
			Expression:  a.attribute.Expr,
			EvalContext: ctx,
		}

		return a.val, hcl.Diagnostics{&diag}
	}

	a.Value = val.AsString()
	a.val = val
	a.ctx = ctx

	return a.val, nil
}
