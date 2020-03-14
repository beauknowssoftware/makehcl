package cmd

import "github.com/beauknowssoftware/makehcl/internal/run"

type RunCommand struct {
	Filename string `short:"f" long:"filename"`
	Verbose  bool   `short:"v" long:"verbose"`
	All      bool   `short:"a" long:"all"`
	DryRun   bool   `short:"d" long:"dry-run"`
}

func (c *RunCommand) Execute(args []string) error {
	var o run.DoOptions
	o.Filename = c.Filename
	o.Verbose = c.Verbose
	o.IgnoreLastModified = c.All
	o.DryRun = c.DryRun
	o.Goal = args

	return run.Do(o)
}
