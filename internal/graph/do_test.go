package graph_test

import (
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"unicode"

	"github.com/beauknowssoftware/makehcl/internal/graph"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
)

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

func do(t *testing.T, o graph.DoOptions) *graph.Graph {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("paniced: %v\n%v", r, string(debug.Stack()))
		}
	}()

	d, diag, err := graph.Do(o)
	if err != nil {
		err = errors.Wrap(err, "failed to do")
		t.Fatal(err)
	}

	if diag.HasErrors() {
		err := errors.Wrap(diag, "failed to do")
		t.Fatal(err)
	}

	return d
}

func failableDo(t *testing.T, o graph.DoOptions) (*graph.Graph, error) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("paniced: %v\n%v", r, string(debug.Stack()))
		}
	}()

	d, diag, err := graph.Do(o)
	if diag.HasErrors() {
		err = errors.Wrap(diag, "failed to do")
		t.Fatal(err)
	}

	return d, err
}

func TestDo(t *testing.T) {
	tests := map[string]struct {
		folder             string
		ignoreParserErrors bool
		options            graph.Options
		expected           string
		error              string
	}{
		"single file imports": {
			folder: "testdata/single_file_imports",
			options: graph.Options{
				GraphType: graph.ImportGraph,
			},
			expected: "digraph {n1 [label=\"make.hcl\"];}",
		},
		"missing import": {
			folder:             "testdata/missing_import",
			ignoreParserErrors: true,
			options: graph.Options{
				GraphType: graph.ImportGraph,
			},
			expected: "digraph {" +
				"n1 [color=\"red\",label=\"import.hcl\"];" +
				"n2 [label=\"make.hcl\"];" +
				"n2 -> n1;" +
				"}",
		},
		"multiple file imports": {
			folder: "testdata/multiple_file_imports",
			options: graph.Options{
				GraphType: graph.ImportGraph,
			},
			expected: "digraph {" +
				"n1 [label=\"import.hcl\"];" +
				"n2 [label=\"make.hcl\"];" +
				"n2 -> n1;" +
				"}",
		},
		"diamond imports": {
			folder: "testdata/diamond_imports",
			options: graph.Options{
				GraphType: graph.ImportGraph,
			},
			expected: "digraph {" +
				"n1 [label=\"diamond.hcl\"];" +
				"n2 [label=\"import.hcl\"];" +
				"n3 [label=\"make.hcl\"];" +
				"n4 [label=\"nested.hcl\"];" +
				"n1 -> n4;" +
				"n2 -> n4;" +
				"n3 -> n1;" +
				"n3 -> n2;" +
				"}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			popd := pushd(t, test.folder)
			defer popd()

			var actual string

			o := graph.DoOptions{
				Options:            test.options,
				IgnoreParserErrors: test.ignoreParserErrors,
			}

			if test.error == "" {
				actual = do(t, o).String()
			} else {
				g, err := failableDo(t, o)
				if diff := cmp.Diff(test.error, err.Error()); diff != "" {
					t.Fatalf("error mismatch (-want,+got):\n%s", diff)
				}
				if g != nil {
					actual = g.String()
				}
			}

			ignoreWhitespace := cmpopts.AcyclicTransformer("IgnoreWhitespace", func(str string) string {
				var b strings.Builder
				b.Grow(len(str))
				for _, ch := range str {
					if !unicode.IsSpace(ch) {
						b.WriteRune(ch)
					}
				}
				return b.String()
			})

			if diff := cmp.Diff(test.expected, actual, ignoreWhitespace); diff != "" {
				t.Fatalf("graph mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}
