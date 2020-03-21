package parse2

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type scope interface {
	childContext(*hcl.EvalContext) *hcl.EvalContext
	set(string, cty.Value)
}

type variableScope struct {
	variables map[string]cty.Value
}

func (s variableScope) childContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	ctx = ctx.NewChild()
	ctx.Variables = s.variables

	return ctx
}

func (s *variableScope) set(name string, val cty.Value) {
	if s.variables == nil {
		s.variables = make(map[string]cty.Value)
	}

	s.variables[name] = val
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
	name     string
	scope    scope
	fillable fillable
}

func (a attribute) fill(ctx *hcl.EvalContext) hcl.Diagnostics {
	if a.scope != nil {
		ctx = a.scope.childContext(ctx)
	}

	val, diag := a.fillable.fill(ctx)

	if a.name != "" {
		a.scope.set(a.name, val)
	}

	return diag
}
