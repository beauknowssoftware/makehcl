package functions

import (
	"path/filepath"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var FilenameSpec = function.Spec{
	Params: []function.Parameter{
		function.Parameter{Type: cty.String},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		path := args[0].AsString()
		ext := filepath.Ext(path)
		result := strings.TrimSuffix(path, ext)
		return cty.StringVal(result), nil
	},
}
