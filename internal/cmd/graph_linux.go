package cmd

import (
	"fmt"
	"strings"

	"github.com/beauknowssoftware/makehcl/internal/graph"
	"github.com/jessevdk/go-flags"
)

type graphType graph.Type

func (t *graphType) Complete(match string) (result []flags.Completion) {
	var options = []graph.Type{
		graph.ImportGraph,
	}

	for _, o := range options {
		if strings.HasPrefix(string(o), match) {
			result = append(result, flags.Completion{
				Item: string(o),
			})
		}
	}

	return
}

type GraphCommand struct {
	GraphType          graphType      `short:"g" long:"graph-type" required:"true"`
	Filename           flags.Filename `short:"f" long:"filename"`
	IgnoreParserErrors bool           `short:"i" long:"ignore-parser-errors"`
}

func (c *GraphCommand) Execute(_ []string) error {
	var o graph.DoOptions
	o.Filename = string(c.Filename)
	o.IgnoreParserErrors = c.IgnoreParserErrors
	o.GraphType = graph.Type(c.GraphType)

	g, diag, err := graph.Do(o)

	if diag.HasErrors() {
		return diag
	}

	if err != nil {
		return err
	}

	fmt.Println(g)

	return nil
}
