package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/beauknowssoftware/makehcl/internal/graph"
	"github.com/jessevdk/go-flags"
	"github.com/windler/dotgraph/renderer"
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
	Show               bool           `short:"s" long:"show"`
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

	if c.Show {
		return c.showGraph(g)
	}

	fmt.Println(g)

	return nil
}

func (c *GraphCommand) showGraph(g *graph.Graph) (err error) {
	var f *os.File

	f, err = ioutil.TempFile("", "*.png")
	if err != nil {
		return
	}

	defer func() {
		err = os.Remove(f.Name())
	}()

	if err := f.Close(); err != nil {
		return err
	}

	r := renderer.PNGRenderer{
		OutputFile: f.Name(),
	}

	r.Render(g.String())

	cmd := exec.Command("open", f.Name(), "-W")
	if err := cmd.Run(); err != nil {
		return err
	}

	return
}
