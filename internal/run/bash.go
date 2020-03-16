package run

import (
	"fmt"
	"os"
	"os/exec"
)

type bashOpts struct {
	verbose    bool
	env        map[string]string
	globalEnv  map[string]string
	shell      string
	shellFlags *string
}

const (
	defaultShell      = "bash"
	defaultShellFlags = "-c"
)

func bash(cmd string, opts bashOpts) error {
	if opts.verbose {
		fmt.Println(cmd)
	}

	shell := opts.shell
	if shell == "" {
		shell = defaultShell
	}

	shellFlags := opts.shellFlags
	if shellFlags == nil {
		c := defaultShellFlags
		shellFlags = &c
	}

	var c *exec.Cmd
	if *shellFlags == "" {
		c = exec.Command(shell, cmd)
	} else {
		c = exec.Command(shell, *shellFlags, cmd)
	}

	c.Env = os.Environ()

	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	e := fmt.Sprintf("HOME=%v", hd)
	c.Env = append(c.Env, e)

	for k, v := range opts.env {
		e := fmt.Sprintf("%v=%v", k, v)
		c.Env = append(c.Env, e)
	}

	for k, v := range opts.globalEnv {
		e := fmt.Sprintf("%v=%v", k, v)
		c.Env = append(c.Env, e)
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
