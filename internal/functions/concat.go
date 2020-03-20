package functions

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func Concat(args []cty.Value, retType cty.Type) (cty.Value, error) {
	if len(args) == 0 {
		return cty.ListValEmpty(cty.List(cty.String)), nil
	}

	var res []cty.Value

	for _, a := range args {
		switch {
		case a.CanIterateElements():
			res = append(res, a.AsValueSlice()...)
		case a.Type() == cty.String:
			res = append(res, a)
		default:
			return cty.UnknownVal(cty.List(cty.String)), fmt.Errorf("expected list of strings, got %v", a.Type().FriendlyName())
		}
	}

	return cty.ListVal(res), nil
}

var ConcatSpec = function.Spec{
	VarParam: &function.Parameter{Type: cty.DynamicPseudoType},
	Type:     function.StaticReturnType(cty.List(cty.String)),
	Impl:     Concat,
}
