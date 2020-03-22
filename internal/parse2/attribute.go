package parse2

import (
	"fmt"

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
	set          setter
	scope        scope
	fillable     fillable
	name         string
	dependencies []string
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

func getDependencies(local string, expr hcl.Expression) (result []string) {
	for _, v := range expr.Variables() {
		root := v.RootName()
		switch root {
		case "var":
			spl := v.SimpleSplit()
			name := spl.Rel[0].(hcl.TraverseAttr).Name
			result = append(result, fmt.Sprintf("var.%v", name))
		default:
			result = append(result, fmt.Sprintf("%v.%v", local, root))
		}
	}

	return
}

type attributeSorter struct {
	attributes      []attribute
	attributeLookup map[string]attribute
	sorted          []attribute
	visited         map[string]bool
	visiting        map[string]bool
}

func (s *attributeSorter) init() (roots []string) {
	s.visited = make(map[string]bool)
	s.visiting = make(map[string]bool)
	s.attributeLookup = make(map[string]attribute)

	dependents := make(map[string][]string)

	for _, a := range s.attributes {
		fmt.Println(a.name)
		s.attributeLookup[a.name] = a

		for _, dep := range a.dependencies {
			fmt.Printf("%v dep %v\n", a.name, dep)
			dependents[dep] = append(dependents[dep], a.name)
		}
	}

	for _, a := range s.attributes {
		deps := dependents[a.name]
		if len(deps) == 0 {
			roots = append(roots, a.name)
		}
	}

	return
}

func (s *attributeSorter) visit(name string) {
	if s.visited[name] {
		return
	}

	if s.visiting[name] {
		panic(fmt.Sprintf("attribute loopback on %v", name))
	}

	s.visiting[name] = true

	a := s.attributeLookup[name]
	for _, dep := range a.dependencies {
		s.visit(dep)
	}

	s.visiting[name] = false
	s.visited[name] = true
	s.sorted = append(s.sorted, a)
}

func (s *attributeSorter) sort() {
	roots := s.init()

	for _, a := range roots {
		s.visit(a)
	}
}
