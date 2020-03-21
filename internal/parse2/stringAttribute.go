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
}

func newStringAttribute(attr *hcl.Attribute, ctx *hcl.EvalContext) (*StringAttribute, hcl.Diagnostics) {
	val, diag := attr.Expr.Value(ctx)
	if diag.HasErrors() {
		return nil, diag
	}

	t := val.Type()
	if t != cty.String {
		diag := hcl.Diagnostic{
			Summary:     "invalid type",
			Detail:      fmt.Sprintf("expected string, got %v", t.FriendlyName()),
			Severity:    hcl.DiagError,
			Subject:     &attr.Range,
			Expression:  attr.Expr,
			EvalContext: ctx,
		}

		return nil, hcl.Diagnostics{&diag}
	}

	sa := StringAttribute{
		Value:     val.AsString(),
		attribute: attr,
		val:       val,
	}

	return &sa, nil
}
