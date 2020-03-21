package parse2

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type setter func(map[string]cty.Value, cty.Value)

func setDirect(name string) setter {
	return func(vars map[string]cty.Value, val cty.Value) {
		vars[name] = val
	}
}

func setOnObject(object, property string) setter {
	return func(vars map[string]cty.Value, val cty.Value) {
		objVal, hasObject := vars[object]
		if !hasObject {
			objVal = cty.ObjectVal(make(map[string]cty.Value))
			vars[object] = objVal
		}

		obj := objVal.AsValueMap()
		if obj == nil {
			obj = make(map[string]cty.Value)
		}

		obj[property] = val
		vars[object] = cty.ObjectVal(obj)
	}
}

type scope interface {
	childContext(*hcl.EvalContext) *hcl.EvalContext
	set(setter, cty.Value)
}

type variableScope struct {
	variables map[string]cty.Value
}

func (s variableScope) childContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	ctx = ctx.NewChild()
	ctx.Variables = s.variables

	return ctx
}

func (s *variableScope) set(setter setter, val cty.Value) {
	if s.variables == nil {
		s.variables = make(map[string]cty.Value)
	}

	setter(s.variables, val)
}

type nestedScope struct {
	outer scope
	variableScope
}

func (s nestedScope) childContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	ctx = s.outer.childContext(ctx)
	return s.variableScope.childContext(ctx)
}

type attribute struct {
	set      setter
	scope    scope
	fillable fillable
}

func (a attribute) fill(ctx *hcl.EvalContext) hcl.Diagnostics {
	if a.scope != nil {
		ctx = a.scope.childContext(ctx)
	}

	val, diag := a.fillable.fill(ctx)

	if a.set != nil {
		a.scope.set(a.set, val)
	}

	return diag
}
