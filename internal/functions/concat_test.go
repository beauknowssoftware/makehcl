package functions_test

import (
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/functions"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

func concat(t *testing.T, args []cty.Value, retType cty.Type) cty.Value {
	v, err := functions.Concat(args, retType)
	if err != nil {
		err = errors.Wrap(err, "failed to concat")
		t.Fatal(err)
	}

	return v
}

func TestConcat(t *testing.T) {
	tests := map[string]struct {
		args     []cty.Value
		retType  cty.Type
		expected cty.Value
	}{
		"empty": {
			expected: cty.ListValEmpty(cty.List(cty.String)),
		},
		"single list": {
			args: []cty.Value{
				cty.ListVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			}),
		},
		"single string": {
			args: []cty.Value{
				cty.StringVal("a"),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
			}),
		},
		"multiple strings": {
			args: []cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			}),
		},
		"mixed": {
			args: []cty.Value{
				cty.StringVal("a"),
				cty.ListVal([]cty.Value{
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
				cty.StringVal("d"),
				cty.ListVal([]cty.Value{
					cty.StringVal("e"),
				}),
				cty.StringVal("f"),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
				cty.StringVal("d"),
				cty.StringVal("e"),
				cty.StringVal("f"),
			}),
		},
		"multiple lists": {
			args: []cty.Value{
				cty.ListVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
				cty.ListVal([]cty.Value{
					cty.StringVal("d"),
					cty.StringVal("e"),
					cty.StringVal("f"),
				}),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
				cty.StringVal("d"),
				cty.StringVal("e"),
				cty.StringVal("f"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := concat(t, test.args, test.retType)

			if !test.expected.RawEquals(actual) {
				t.Fatalf("\nexpected %v\ngot %v", test.expected.GoString(), actual.GoString())
			}
		})
	}
}
