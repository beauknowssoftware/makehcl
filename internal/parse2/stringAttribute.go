package parse2

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type StringAttribute struct {
	Value     string
	attribute *hcl.Attribute
	val       cty.Value
	ctx       *hcl.EvalContext
}

func newStringAttribute(attr *hcl.Attribute, ctx *hcl.EvalContext) (sa *StringAttribute, diag hcl.Diagnostics) {
	sa = &StringAttribute{
		attribute: attr,
	}
	diag = sa.fill(ctx)

	return
}

func (a *StringAttribute) fill(ctx *hcl.EvalContext) hcl.Diagnostics {
	val, diag := a.attribute.Expr.Value(ctx)
	if diag.HasErrors() {
		return diag
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

		return hcl.Diagnostics{&diag}
	}

	a.Value = val.AsString()
	a.val = val
	a.ctx = ctx

	return nil
}
