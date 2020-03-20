package cmd

import (
	"github.com/beauknowssoftware/makehcl/internal/run"
	"github.com/jessevdk/go-flags"
)

type RunCommand struct {
	Filename flags.Filename `short:"f" long:"filename"`
	Verbose  bool           `short:"v" long:"verbose"`
	All      bool           `short:"a" long:"all"`
	DryRun   bool           `short:"d" long:"dry-run"`
	Goal     Goal           `positional-args:"yes"`
}

func (c *RunCommand) Execute(_ []string) error {
	var o run.DoOptions
	o.Filename = string(c.Filename)
	o.Verbose = c.Verbose
	o.IgnoreLastModified = c.All
	o.DryRun = c.DryRun
	o.Goal = c.Goal.strings()

	return run.Do(o)
}
