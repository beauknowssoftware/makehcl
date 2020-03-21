package graph

import (
	"sort"

	"github.com/beauknowssoftware/makehcl/internal/parse2"
	"github.com/emicklei/dot"
)

func constructImportGraph(d *parse2.Definition) *Graph {
	g := dot.NewGraph(dot.Directed)

	nodeMap := make(map[string]dot.Node)

	filenames := make([]string, 0, len(d.Files))
	for _, f := range d.Files {
		filenames = append(filenames, f.Name)
	}

	sort.Strings(filenames)

	for _, name := range filenames {
		n := g.Node(name)

		f := d.Files[name]
		if !f.HasContents() {
			n.Attr("color", "red")
		}

		nodeMap[name] = n
	}

	for _, name := range filenames {
		n1 := nodeMap[name]

		f := d.Files[name]
		imports := make([]string, 0, len(f.ImportBlocks))

		for _, imp := range f.ImportBlocks {
			if imp.File == nil {
				continue
			}

			imports = append(imports, imp.File.Value)
		}

		sort.Strings(imports)

		for _, imp := range f.ImportBlocks {
			if imp.File == nil {
				continue
			}

			n2 := nodeMap[imp.File.Value]
			g.Edge(n1, n2)
		}
	}

	return g
}
