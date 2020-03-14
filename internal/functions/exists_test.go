package functions_test

import (
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/functions"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

func exists(t *testing.T, args []cty.Value, retType cty.Type) cty.Value {
	v, err := functions.Exists(args, retType)
	if err != nil {
		err = errors.Wrap(err, "failed to exists")
		t.Fatal(err)
	}
	return v
}

func TestExists(t *testing.T) {
	tests := map[string]struct {
		args     []cty.Value
		retType  cty.Type
		expected cty.Value
	}{
		"object exists": {
			args: []cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"test": cty.StringVal("value"),
				}),
				cty.StringVal("test"),
			},
			expected: cty.BoolVal(true),
		},
		"object not exists": {
			args: []cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"test": cty.StringVal("value"),
				}),
				cty.StringVal("another"),
			},
			expected: cty.BoolVal(false),
		},
		"map exists": {
			args: []cty.Value{
				cty.MapVal(map[string]cty.Value{
					"test": cty.StringVal("value"),
				}),
				cty.StringVal("test"),
			},
			expected: cty.BoolVal(true),
		},
		"map not exists": {
			args: []cty.Value{
				cty.MapVal(map[string]cty.Value{
					"test": cty.StringVal("value"),
				}),
				cty.StringVal("another"),
			},
			expected: cty.BoolVal(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := exists(t, test.args, test.retType)

			if !test.expected.RawEquals(actual) {
				t.Fatalf("\nexpected %v\ngot %v", test.expected.GoString(), actual.GoString())
			}
		})
	}
}
