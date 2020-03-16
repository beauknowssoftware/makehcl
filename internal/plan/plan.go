package plan

import (
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

func visit(o Options, d definition.Definition, t definition.Target, p *Plan, visited map[definition.Target]bool, dt map[definition.Target]*time.Time) error {
	if visited[t] {
		return nil
	}

	mt, err := file.ModTime(t)
	if err != nil {
		return err
	}

	dt[t] = mt

	r := d.Rule(t)
	if r == nil && mt != nil {
		return nil
	}

	shouldVisit := o.IgnoreLastModified || mt == nil // mt should only be null if file does not exist

	for _, dep := range r.Dependencies {
		err := visit(o, d, dep, p, visited, dt)
		if err != nil {
			return err
		}

		shouldVisit = shouldVisit || visited[dep] || (mt != nil && dt[dep].After(*mt))
	}

	if shouldVisit {
		*p = append(*p, t)
		visited[t] = true
	}

	return nil
}

func Definition(d definition.Definition, o Options) (Plan, error) {
	g := d.EffectiveGoal(o.Goal)

	visited := make(map[definition.Target]bool)
	dt := make(map[definition.Target]*time.Time)
	p := Plan{}

	for _, t := range g {
		if err := visit(o, d, t, &p, visited, dt); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func Do(o DoOptions) (Plan, error) {
	d, err := parse.Parse(o.ParseOptions)
	if err != nil {
		return nil, err
	}

	return Definition(*d, o.Options)
}
