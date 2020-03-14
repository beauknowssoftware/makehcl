package functions_test

import (
	"os"
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/functions"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

func glob(t *testing.T, args []cty.Value, retType cty.Type) cty.Value {
	v, err := functions.Glob(args, retType)
	if err != nil {
		err = errors.Wrap(err, "failed to glob")
		t.Fatal(err)
	}
	return v
}

func pushd(t *testing.T, dir string) func() {
	originalDir, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get working directory")
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		err = errors.Wrapf(err, "failed to pushd to %v", dir)
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(originalDir); err != nil {
			err = errors.Wrapf(err, "failed to popd directories to %v", originalDir)
			t.Fatal(err)
		}
	}
}

func TestGlob(t *testing.T) {
	tests := map[string]struct {
		directory string
		args      []cty.Value
		retType   cty.Type
		expected  cty.Value
	}{
		"glob": {
			directory: "testdata/glob",
			args: []cty.Value{
				cty.StringVal("*.txt"),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("test.txt"),
				cty.StringVal("test2.txt"),
			}),
		},
		"nested glob": {
			directory: "testdata/nested_glob",
			args: []cty.Value{
				cty.StringVal("**.txt"),
			},
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("nested/test.txt"),
				cty.StringVal("nested/test2.txt"),
				cty.StringVal("test.txt"),
				cty.StringVal("test2.txt"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			popd := pushd(t, test.directory)
			defer popd()

			actual := glob(t, test.args, test.retType)

			if !test.expected.RawEquals(actual) {
				t.Fatalf("\nexpected %v\ngot %v", test.expected.GoString(), actual.GoString())
			}
		})
	}
}
