package graph

import (
	"github.com/beauknowssoftware/makehcl/internal/parse2"
	"github.com/hashicorp/hcl/v2"
)

func Do(o DoOptions) (*Graph, hcl.Diagnostics, error) {
	po := parse2.Options{
		Filename:       o.Filename,
		StopAfterStage: parse2.StopAfterImports,
	}

	d, diag := parse2.Do(po)
	if diag.HasErrors() && !o.IgnoreParserErrors {
		return nil, diag, nil
	}

	g, err := ConstructGraph(&d, o.Options)

	return g, nil, err
}
