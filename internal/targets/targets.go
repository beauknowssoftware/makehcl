package targets

import (
	"sort"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/beauknowssoftware/makehcl/internal/parse"
)

type Plan []definition.Target

type ParseOptions = parse.Options

type Options struct {
	Sort bool
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
		r = append(r, rl.Target)
	}

	if o.Sort {
		sort.Strings(r)
	}

	return r, nil
}
