package cmd

import (
	"fmt"

	"github.com/beauknowssoftware/makehcl/internal/plan"
)

type PlanCommand struct {
	Filename string `short:"f" long:"filename"`
	All      bool   `short:"a" long:"all"`
}

func (c *PlanCommand) Execute(args []string) error {
	var o plan.DoOptions
	o.Filename = c.Filename
	o.IgnoreLastModified = c.All
	o.Goal = args

	p, err := plan.Do(o)
	if err != nil {
		return err
	}

	for _, e := range p {
		fmt.Println(e)
	}

	return nil
}
