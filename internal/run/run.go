package run

import (
	"fmt"

	"github.com/beauknowssoftware/makehcl/internal/parse"
	"github.com/beauknowssoftware/makehcl/internal/plan"
)

type ParseOptions = parse.Options
type PlanOptions = plan.Options

type Options struct {
	Verbose bool
	DryRun  bool
}

type DoOptions struct {
	ParseOptions
	PlanOptions
	Options
}

func Do(o DoOptions) error {
	d, err := parse.Parse(o.ParseOptions)
	if err != nil {
		return err
	}

	p, err := plan.Definition(*d, o.PlanOptions)
	if err != nil {
		return err
	}

	for _, t := range p {
		r := d.Rule(t)

		if o.DryRun {
			fmt.Println(r.Command)
		} else {
			opts := bashOpts{
				verbose:    o.Verbose,
				env:        r.Environment,
				globalEnv:  d.GlobalEnvironment,
				shell:      d.Shell,
				shellFlags: d.ShellFlags,
			}
			if r.TeeTarget {
				opts.teeTarget = &t
			}
			if err := bash(r.Command, opts); err != nil {
				return err
			}
		}
	}

	return nil
}
