package cmd

import "github.com/jessevdk/go-flags"

func Exec() error {
	p := flags.NewParser(nil, flags.Default)

	var plan PlanCommand
	_, err := p.AddCommand("plan",
		"Plan execution",
		"The plan command plans execution",
		&plan)

	if err != nil {
		return err
	}

	var run RunCommand
	_, err = p.AddCommand("run",
		"Execute",
		"The run command executes",
		&run)

	if err != nil {
		return err
	}

	var targets TargetsCommand
	_, err = p.AddCommand("targets",
		"Display list of targets",
		"The targets command displays a list of executable targets",
		&targets)

	if err != nil {
		return err
	}

	var completion CompletionCommand
	_, err = p.AddCommand("completion",
		"Show completion script",
		"The completion command prints a script that can be used for shell completion",
		&completion)

	if err != nil {
		return err
	}

	var graph GraphCommand
	_, err = p.AddCommand("graph",
		"Print a graph",
		"The graph command parses the make hcl file and prints a dot formatted graph",
		&graph)

	if err != nil {
		return err
	}

	_, err = p.Parse()

	return err
}
