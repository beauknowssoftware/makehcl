package functions

import (
	"path/filepath"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var PathSpec = function.Spec{
	VarParam: &function.Parameter{Type: cty.String},
	Type:     function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		sa := make([]string, 0, len(args))
		for _, a := range args {
			sa = append(sa, a.AsString())
		}
		result := filepath.Join(sa...)
		return cty.StringVal(result), nil
	},
}
