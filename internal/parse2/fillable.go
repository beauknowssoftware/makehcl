package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

type fillable interface {
	fill(*hcl.EvalContext) hcl.Diagnostics
}
