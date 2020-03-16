package plan

import (
	"fmt"
	"time"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/beauknowssoftware/makehcl/internal/file"
	"github.com/beauknowssoftware/makehcl/internal/parse"
)

type Plan []definition.Target

type Options struct {
	Goal               definition.Goal
	IgnoreLastModified bool
}

type ParseOptions = parse.Options

type DoOptions struct {
	ParseOptions
	Options
}

type planVisitor struct {
	o       Options
	d       definition.Definition
	p       Plan
	visited map[definition.Target]bool
	dt      map[definition.Target]*time.Time
}

func (v *planVisitor) visit(t definition.Target) error {
	if v.visited[t] {
		return nil
	}

	mt, err := file.ModTime(t)
	if err != nil {
		return err
	}

	v.dt[t] = mt

	r := v.d.Rule(t)
	if r == nil && mt != nil {
		return nil
	}

	if r == nil {
		return fmt.Errorf("unknown target %v", t)
	}

	shouldVisit := v.o.IgnoreLastModified || mt == nil // mt should only be null if file does not exist

	for _, dep := range r.Dependencies {
		err := v.visit(dep)
		if err != nil {
			return err
		}

		shouldVisit = shouldVisit || v.visited[dep] || (mt != nil && v.dt[dep].After(*mt))
	}

	if shouldVisit {
		v.p = append(v.p, t)
		v.visited[t] = true
	}

	return nil
}

func Definition(d definition.Definition, o Options) (Plan, error) {
	g := d.EffectiveGoal(o.Goal)

	var visitor planVisitor
	visitor.visited = make(map[definition.Target]bool)
	visitor.dt = make(map[definition.Target]*time.Time)
	visitor.d = d
	visitor.o = o

	for _, t := range g {
		if err := visitor.visit(t); err != nil {
			return nil, err
		}
	}

	return visitor.p, nil
}

func Do(o DoOptions) (Plan, error) {
	d, err := parse.Parse(o.ParseOptions)
	if err != nil {
		return nil, err
	}

	return Definition(*d, o.Options)
}
