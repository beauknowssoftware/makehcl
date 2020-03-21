package parse2

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type generic struct {
	attribute *hcl.Attribute
	val       cty.Value
	ctx       *hcl.EvalContext
}

func (a *generic) fill(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	val, diag := a.attribute.Expr.Value(ctx)
	if diag.HasErrors() {
		return a.val, diag
	}

	a.val = val
	a.ctx = ctx

	return val, nil
}
