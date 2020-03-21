package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

type attribute interface {
	fill(*hcl.EvalContext) hcl.Diagnostics
}
