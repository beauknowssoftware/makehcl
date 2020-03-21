package graph

import (
	"github.com/beauknowssoftware/makehcl/internal/parse2"
	"github.com/emicklei/dot"
)

func constructImportGraph(d *parse2.Definition) *Graph {
	g := dot.NewGraph(dot.Directed)

	for _, f := range d.Files {
		g.Node(f.Name)
	}

	return g
}
