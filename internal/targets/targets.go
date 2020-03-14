package targets

import (
	"sort"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/beauknowssoftware/makehcl/internal/parse"
)

type Plan []definition.Target

type ParseOptions = parse.Options

type Options struct {
	Sort        bool
	CommandOnly bool
	RuleOnly    bool
}

type DoOptions struct {
	ParseOptions
	Options
}

func Do(o DoOptions) ([]definition.Target, error) {
	d, err := parse.Parse(o.ParseOptions)
	if err != nil {
		return nil, err
	}

	rules := d.Rules()
	r := make([]definition.Target, 0, len(rules))
	for _, rl := range rules {
		if o.RuleOnly && !rl.IsPhony {
			r = append(r, rl.Target)
		} else if o.CommandOnly && rl.IsPhony {
			r = append(r, rl.Target)
		} else if !o.RuleOnly && !o.CommandOnly {
			r = append(r, rl.Target)
		}
	}

	if o.Sort {
		sort.Strings(r)
	}

	return r, nil
}
