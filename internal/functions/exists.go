package functions

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func Exists(args []cty.Value, _ cty.Type) (cty.Value, error) {
	key := args[1].AsString()

	t := args[0].Type()

	switch {
	case t.IsObjectType():
		res := t.HasAttribute(key)
		return cty.BoolVal(res), nil
	case t.IsMapType():
		m := args[0].AsValueMap()
		_, res := m[key]

		return cty.BoolVal(res), nil
	default:
		return cty.False, fmt.Errorf("expected map or object got %v", t.FriendlyName())
	}
}

var ExistsSpec = function.Spec{
	Params: []function.Parameter{
		function.Parameter{Type: cty.DynamicPseudoType},
		function.Parameter{Type: cty.String},
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: Exists,
}
