package parse2_test

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/parse2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
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

func do(t *testing.T, o parse2.Options) parse2.Definition {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("paniced: %v\n%v", r, string(debug.Stack()))
		}
	}()

	d, diag := parse2.Do(o)
	if diag.HasErrors() {
		err := errors.Wrap(diag, "failed to do")
		t.Fatal(err)
	}

	return d
}

func TestDo(t *testing.T) {
	tests := map[string]struct {
		folder     string
		options    parse2.Options
		definition parse2.Definition
		error      string
	}{
		"filename": {
			folder: "testdata/filename",
			options: parse2.Options{
				Filename: "filename.hcl",
			},
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"filename.hcl": {
						Name: "filename.hcl",
					},
				},
			},
		},
		"missing file": {
			folder: "testdata/missing_file",
			error:  "<nil>: Failed to read file; The configuration file \"make.hcl\" could not be read.",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
					},
				},
			},
		},
		"missing import": {
			folder: "testdata/missing_import",
			error:  "<nil>: Failed to read file; The configuration file \"import.hcl\" could not be read.",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
					},
				},
			},
		},
		"empty": {
			folder: "testdata/empty",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
					},
				},
			},
		},
		"import": {
			folder: "testdata/import",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
					},
				},
			},
		},
		"import only": {
			folder: "testdata/import_only",
			options: parse2.Options{
				StopAfterStage: parse2.StopAfterImports,
			},
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "make.hcl",
								},
							},
						},
					},
				},
			},
		},
		"import loop": {
			folder: "testdata/import_loop",
			error:  "import.hcl:2,3-20: Import cycle detected; Cycle occurred when importing make.hcl",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "make.hcl",
								},
							},
						},
					},
				},
			},
		},
		"nested import": {
			folder: "testdata/nested_import",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "nested.hcl",
								},
							},
						},
					},
					"nested.hcl": {
						Name: "nested.hcl",
					},
				},
			},
		},
		"diamond import": {
			folder: "testdata/diamond_import",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "import.hcl",
								},
							},
							{
								File: &parse2.StringAttribute{
									Value: "diamond.hcl",
								},
							},
						},
					},
					"diamond.hcl": {
						Name: "diamond.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "nested.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.StringAttribute{
									Value: "nested.hcl",
								},
							},
						},
					},
					"nested.hcl": {
						Name: "nested.hcl",
					},
				},
			},
		},
	}

	ignoreUnexported := cmpopts.IgnoreUnexported(
		parse2.Definition{},
		parse2.File{},
		parse2.ImportBlock{},
		parse2.StringAttribute{},
	)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			popd := pushd(t, test.folder)
			defer popd()

			var definition parse2.Definition

			if test.error == "" {
				definition = do(t, test.options)
			} else {
				var diag hcl.Diagnostics
				definition, diag = parse2.Do(test.options)
				if diff := cmp.Diff(test.error, diag.Error()); diff != "" {
					t.Fatalf("error mismatch (-want,+got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(test.definition, definition, ignoreUnexported); diff != "" {
				t.Fatalf("definition mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}
