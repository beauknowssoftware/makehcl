package parse2

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type fillable interface {
	fill(*hcl.EvalContext) (cty.Value, hcl.Diagnostics)
}
