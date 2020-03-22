package parse2_test

import (
	"fmt"
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
								File: &parse2.String{
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
		"var": {
			folder: "testdata/var",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						RuleBlocks: []*parse2.RuleBlock{
							{
								Target: &parse2.String{
									Value: "test.txt",
								},
								Command: &parse2.StringArray{
									Value: []string{"touch test.txt"},
								},
							},
						},
					},
				},
			},
		},
		"rule target var": {
			folder: "testdata/rule_target_var",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						RuleBlocks: []*parse2.RuleBlock{
							{
								Target: &parse2.String{
									Value: "test.txt",
								},
								Command: &parse2.StringArray{
									Value: []string{"touch test.txt"},
								},
							},
						},
					},
				},
			},
		},
		"commands": {
			folder: "testdata/commands",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						CommandBlocks: []*parse2.CommandBlock{
							{
								Name: "test",
								Command: &parse2.StringArray{
									Value: []string{"touch test.txt"},
								},
							},
						},
					},
				},
			},
		},
		"default goal": {
			folder: "testdata/default_goal",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						DefaultGoal: &parse2.StringArray{
							Value: []string{"test.txt"},
						},
						Name: "make.hcl",
					},
				},
			},
		},
		"multiple default goals": {
			folder: "testdata/multiple_default_goals",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						DefaultGoal: &parse2.StringArray{
							Value: []string{"test.txt", "test2.txt"},
						},
						Name: "make.hcl",
					},
				},
			},
		},
		"rules": {
			folder: "testdata/rules",
			definition: parse2.Definition{
				Files: map[string]*parse2.File{
					"make.hcl": {
						Name: "make.hcl",
						RuleBlocks: []*parse2.RuleBlock{
							{
								Target: &parse2.String{
									Value: "test.txt",
								},
								Command: &parse2.StringArray{
									Value: []string{"touch test.txt"},
								},
							},
						},
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
								File: &parse2.String{
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
								File: &parse2.String{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.String{
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
								File: &parse2.String{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.String{
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
								File: &parse2.String{
									Value: "import.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.String{
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
								File: &parse2.String{
									Value: "import.hcl",
								},
							},
							{
								File: &parse2.String{
									Value: "diamond.hcl",
								},
							},
						},
					},
					"diamond.hcl": {
						Name: "diamond.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.String{
									Value: "nested.hcl",
								},
							},
						},
					},
					"import.hcl": {
						Name: "import.hcl",
						ImportBlocks: []*parse2.ImportBlock{
							{
								File: &parse2.String{
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
		parse2.RuleBlock{},
		parse2.CommandBlock{},
		parse2.String{},
		parse2.StringArray{},
	)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fmt.Println(name)

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
