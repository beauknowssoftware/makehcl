package cmd

import (
	"errors"
	"fmt"

	"github.com/beauknowssoftware/makehcl/internal/targets"
	"github.com/jessevdk/go-flags"
)

type TargetsCommand struct {
	Filename    flags.Filename `short:"f" long:"filename"`
	Sort        bool           `short:"s" long:"sort"`
	RuleOnly    bool           `short:"r" long:"rule-only"`
	CommandOnly bool           `short:"c" long:"command-only"`
}

func (c *TargetsCommand) Execute(_ []string) error {
	if c.RuleOnly && c.CommandOnly {
		return errors.New("cannot specify rule only and command only at the same time")
	}

	var o targets.DoOptions
	o.Filename = string(c.Filename)
	o.Sort = c.Sort
	o.RuleOnly = c.RuleOnly
	o.CommandOnly = c.CommandOnly

	ts, err := targets.Do(o)
	if err != nil {
		return err
	}

	for _, t := range ts {
		fmt.Println(t)
	}

	return nil
}
