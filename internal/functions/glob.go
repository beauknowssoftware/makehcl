package functions

import (
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func Glob(args []cty.Value, _ cty.Type) (cty.Value, error) {
	globString := args[0].AsString()
	g, err := glob.Compile(globString, filepath.Separator)

	if err != nil {
		err = errors.Wrapf(err, "invalid glob %v", globString)
		return cty.UnknownVal(cty.String), err
	}

	var vals []cty.Value

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if g.Match(path) {
			vals = append(vals, cty.StringVal(path))
		}
		return nil
	})

	if err != nil {
		return cty.UnknownVal(cty.String), err
	}

	if len(vals) == 0 {
		return cty.ListValEmpty(cty.String), nil
	}

	return cty.ListVal(vals), nil
}

var GlobSpec = function.Spec{
	Params: []function.Parameter{
		function.Parameter{Type: cty.String},
	},
	Type: function.StaticReturnType(cty.List(cty.String)),
	Impl: Glob,
}
