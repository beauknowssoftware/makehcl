package graph

import (
	"errors"

	"github.com/beauknowssoftware/makehcl/internal/parse2"
)

func ConstructGraph(d *parse2.Definition, o Options) (*Graph, error) {
	if o.GraphType != ImportGraph {
		return nil, errors.New("invalid graph type")
	}

	return constructImportGraph(d), nil
}
