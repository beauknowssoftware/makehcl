package cmd

import (
	"fmt"

	"github.com/beauknowssoftware/makehcl/internal/plan"
	"github.com/jessevdk/go-flags"
)

type PlanCommand struct {
	Filename flags.Filename `short:"f" long:"filename"`
	All      bool           `short:"a" long:"all"`
	Goal     Goal           `positional-args:"yes"`
}

func (c *PlanCommand) Execute(_ []string) error {
	var o plan.DoOptions
	o.Filename = string(c.Filename)
	o.IgnoreLastModified = c.All
	o.Goal = c.Goal.strings()

	p, err := plan.Do(o)
	if err != nil {
		return err
	}

	for _, e := range p {
		fmt.Println(e)
	}

	return nil
}
