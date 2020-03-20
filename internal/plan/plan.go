package plan

import (
	"fmt"
	"time"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/beauknowssoftware/makehcl/internal/file"
	"github.com/beauknowssoftware/makehcl/internal/parse"
)

const (
	DoesNotExist        = ReasonType("does not exist")
	OlderThanDependency = ReasonType("older than dependency")
	DependencyPlanned   = ReasonType("dependency planned")
)

type ReasonType string

type Reason struct {
	Target     definition.Target
	ReasonType ReasonType
}

type Item struct {
	Target definition.Target
	Reason Reason
}
type Plan []Item

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

func (v planVisitor) getTargetInfo(t definition.Target) (*definition.Rule, *time.Time, error) {
	mt, err := file.ModTime(t)
	if err != nil {
		return nil, nil, err
	}

	v.dt[t] = mt

	r := v.d.Rule(t)
	// if rule does not exist but there is a last modified time
	// then we are looking at a pre-existing file dependency
	if r == nil && mt != nil {
		return nil, mt, nil
	}

	// if the rule does not exist, and there is no last modified time
	// then we're missing a target
	if r == nil {
		return nil, nil, fmt.Errorf("unknown target %v", t)
	}

	return r, mt, nil
}

func (v *planVisitor) visitDependency(dep definition.Target, mt *time.Time) (*Reason, error) {
	err := v.visit(dep)
	if err != nil {
		return nil, err
	}

	if v.visited[dep] {
		return &Reason{dep, DependencyPlanned}, nil
	} else if mt != nil && v.dt[dep].After(*mt) {
		return &Reason{dep, OlderThanDependency}, nil
	}

	return nil, nil
}

func (v *planVisitor) visit(t definition.Target) error {
	if v.visited[t] {
		return nil
	}

	r, mt, err := v.getTargetInfo(t)
	if err != nil {
		return err
	}

	if r == nil {
		return nil
	}

	var reason *Reason

	// mt should only be null if file does not exist
	if mt == nil {
		reason = &Reason{t, DoesNotExist}
	}

	for _, dep := range r.Dependencies {
		depReason, err := v.visitDependency(dep, mt)
		if err != nil {
			return err
		}

		if reason == nil {
			reason = depReason
		}
	}

	v.finalizeVisit(t, reason)

	return nil
}

func (v *planVisitor) finalizeVisit(t definition.Target, reason *Reason) {
	if reason != nil || v.o.IgnoreLastModified {
		i := Item{Target: t}
		if reason != nil && !v.o.IgnoreLastModified {
			i.Reason = *reason
		}

		v.p = append(v.p, i)
		v.visited[t] = true
	}
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
