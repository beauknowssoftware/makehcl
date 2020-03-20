package cmd

import (
	"fmt"
	"os"
)

type CompletionCommand struct {
}

func (c *CompletionCommand) Execute(_ []string) error {
	name := os.Args[0]
	fmt.Printf(`_%v() {
    # All arguments except the first one
    args=("${COMP_WORDS[@]:1:$COMP_CWORD}")

    # Only split on newlines
    local IFS=$'\n'

    # Call completion (note that the first element of COMP_WORDS is
    # the executable itself)
    COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
    return 0
}

complete -F _%v %v`+"\n", name, name, name)

	return nil
}
