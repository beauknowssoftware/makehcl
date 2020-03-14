package parse

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type global struct {
	globalType string
	attr       *hcl.Attribute
}

func (g global) name() globalName {
	return globalName{
		globalType: g.globalType,
		name:       g.attr.Name,
	}
}

type globalName struct {
	globalType string
	name       string
}

type globalSorter struct {
	globals  map[globalName]global
	sorted   []global
	visited  map[globalName]bool
	visiting map[globalName]bool
}

func (s *globalSorter) visit(g global) {
	if s.visited[g.name()] {
		return
	}

	if s.visiting[g.name()] {
		panic(fmt.Sprintf("env loopback on %v", g.name()))
	}

	s.visiting[g.name()] = true

	for _, v := range g.attr.Expr.Variables() {
		globalType := v.RootName()
		spl := v.SimpleSplit()
		if spl.Rel == nil || len(spl.Rel) == 0 {
			continue
		}
		name := spl.Rel[0].(hcl.TraverseAttr).Name
		gName := globalName{globalType, name}
		if g.name() == gName && globalType == "env" {
			continue
		}
		if g, hasGlobal := s.globals[gName]; hasGlobal {
			s.visit(g)
		}
	}

	s.visiting[g.name()] = false
	s.visited[g.name()] = true
	s.sorted = append(s.sorted, g)
}

func (s *globalSorter) sort() {
	s.visited = make(map[globalName]bool)
	s.visiting = make(map[globalName]bool)
	for _, a := range s.globals {
		s.visit(a)
	}
}

func fillGlobals(attrSets map[string]map[string]*hcl.Attribute, ctx *hcl.EvalContext) (map[string]string, error) {
	s := globalSorter{
		globals: make(map[globalName]global),
	}
	for globalType, attr := range attrSets {
		for _, a := range attr {
			g := global{
				globalType: globalType,
				attr:       a,
			}
			s.globals[g.name()] = g
		}
	}
	s.sort()
	vars := make(map[string]cty.Value)
	envResult := make(map[string]string)
	envs := ctx.Variables["env"].AsValueMap()

	for _, a := range s.sorted {
		if a.globalType == "var" {
			val, diag := a.attr.Expr.Value(ctx)
			vars[a.attr.Name] = val
			if diag.HasErrors() {
				return nil, errors.Wrap(diag, "failed to get var attributes")
			}
			ctx.Variables["var"] = cty.ObjectVal(vars)
		} else if a.globalType == "env" {
			val, diag := a.attr.Expr.Value(ctx)
			envs[a.attr.Name] = val
			envResult[a.attr.Name] = val.AsString()
			if diag.HasErrors() {
				return nil, errors.Wrap(diag, "failed to get env attributes")
			}
			ctx.Variables["env"] = cty.ObjectVal(envs)
		}
	}

	return envResult, nil
}
