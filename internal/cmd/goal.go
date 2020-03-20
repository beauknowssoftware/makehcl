package cmd

import (
	"strings"

	"github.com/beauknowssoftware/makehcl/internal/targets"
	"github.com/jessevdk/go-flags"
)

type Target string

func (t *Target) Complete(match string) []flags.Completion {
	var o targets.DoOptions
	o.Sort = true

	ts, err := targets.Do(o)
	if err != nil {
		return []flags.Completion{}
	}

	res := make([]flags.Completion, 0, len(ts))

	for _, t := range ts {
		if strings.HasPrefix(t, match) {
			res = append(res, flags.Completion{
				Item: t,
			})
		}
	}

	return res
}

type Goal struct {
	Targets []Target
}

func (g Goal) strings() []string {
	res := make([]string, len(g.Targets))
	for i, t := range g.Targets {
		res[i] = string(t)
	}

	return res
}
