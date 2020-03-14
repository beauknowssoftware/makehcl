package cmd

import (
	"fmt"

	"github.com/beauknowssoftware/makehcl/internal/targets"
)

type TargetsCommand struct {
	Filename string `short:"f" long:"filename"`
	Sort     bool   `short:"s" long:"sort"`
}

func (c *TargetsCommand) Execute(_ []string) error {
	var o targets.DoOptions
	o.Filename = c.Filename
	o.Sort = c.Sort

	ts, err := targets.Do(o)
	if err != nil {
		return err
	}

	for _, t := range ts {
		fmt.Println(t)
	}
	return nil
}
