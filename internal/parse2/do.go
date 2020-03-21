package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

func Do(o Options) (Definition, hcl.Diagnostics) {
	var p Parser
	p.Options = o

	diag := p.Parse()

	return p.Definition, diag
}
