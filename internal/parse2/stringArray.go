package parse2

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type StringArray struct {
	Value     []string
	attribute *hcl.Attribute
	val       cty.Value
	ctx       *hcl.EvalContext
}

func isStringType(_ cty.Value, val cty.Value) (stop bool) {
	return val.Type() == cty.String
}

func (a *StringArray) fill(ctx *hcl.EvalContext) hcl.Diagnostics {
	val, diag := a.attribute.Expr.Value(ctx)
	if diag.HasErrors() {
		return diag
	}

	t := val.Type()

	if t == cty.String {
		v := val.AsString()
		a.Value = []string{v}
		a.val = val
		a.ctx = ctx

		return nil
	}

	if !val.CanIterateElements() || val.ForEachElement(isStringType) {
		diag := hcl.Diagnostic{
			Summary:     "invalid type",
			Detail:      fmt.Sprintf("expected iterable strings, got %v", t.FriendlyName()),
			Severity:    hcl.DiagError,
			Subject:     &a.attribute.Range,
			Expression:  a.attribute.Expr,
			EvalContext: ctx,
		}

		return hcl.Diagnostics{&diag}
	}

	slc := val.AsValueSlice()
	a.Value = make([]string, 0, len(slc))
	a.val = val
	a.ctx = ctx

	for _, i := range slc {
		a.Value = append(a.Value, i.AsString())
	}

	return nil
}
