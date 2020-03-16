package plan_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/plan"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

var (
	local = false
)

type filename = string
type fileContents = []byte

func createFile(t *testing.T, f filename, data fileContents) {
	dir := filepath.Dir(f)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		err = errors.Wrapf(err, "failed to create directory %v", dir)
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(f, data, os.ModePerm); err != nil {
		err = errors.Wrapf(err, "failed to write file %v", f)
		t.Fatal(err)
	}
}

func copyFile(t *testing.T, src, dest string) {
	out, err := os.Create(dest)
	if err != nil {
		err = errors.Wrapf(err, "failed to create %v", dest)
		t.Fatal(err)
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		err = errors.Wrapf(err, "failed to open %v", src)
		t.Fatal(err)
	}
	defer in.Close()

	if _, err := io.Copy(out, in); err != nil {
		err = errors.Wrapf(err, "failed to copy from %v to %v", src, dest)
		t.Fatal(err)
	}
}

func safeDo(t *testing.T, o plan.DoOptions) plan.Plan {
	p, err := plan.Do(o)
	if err != nil {
		err = errors.Wrap(err, "failed to plan")
		t.Fatal(err)
	}

	return p
}

func copyDir(t *testing.T, src, dest string) {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		err = errors.Wrapf(err, "failed to list %v", src)
		t.Fatal(err)
	}

	for _, entry := range entries {
		srcFile := filepath.Join(src, entry.Name())
		destFile := filepath.Join(dest, entry.Name())
		copyFile(t, srcFile, destFile)
	}
}

func copyToTemp(t *testing.T, src string) string {
	tempDir, err := ioutil.TempDir("", "makehcl")
	if err != nil {
		err = errors.Wrap(err, "failed to create temporary directory")
		t.Fatal(err)
	}

	copyDir(t, src, tempDir)

	return tempDir
}

func pushd(t *testing.T, dir string) func() {
	originalDir, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get working directory")
		t.Fatal(err)
	}

	if !local {
		dir = copyToTemp(t, dir)
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

		if !local {
			if err := os.RemoveAll(dir); err != nil {
				err = errors.Wrapf(err, "failed to remove %v", dir)
				t.Fatal(err)
			}
		}
	}
}

func TestRun(t *testing.T) {
	tests := map[string]struct {
		folder   string
		existing []map[filename]fileContents
		expected plan.Plan
		options  plan.DoOptions
	}{
		"simple": {
			folder: "testdata/simple",
			expected: plan.Plan{
				"test.txt",
			},
		},
		"dependencies": {
			folder: "testdata/dependencies",
			expected: plan.Plan{
				"test2.txt",
				"test.txt",
				"test3.txt",
			},
		},
		"existing files": {
			folder: "testdata/dependencies",
			existing: []map[filename]fileContents{
				{"test2.txt": {}},
			},
			expected: plan.Plan{
				"test.txt",
				"test3.txt",
			},
		},
		"ignore last modified": {
			folder: "testdata/ignore_last_modified",
			options: plan.DoOptions{
				Options: plan.Options{
					IgnoreLastModified: true,
				},
			},
			existing: []map[filename]fileContents{
				{"test2.txt": {}},
				{"test.txt": {}},
				{"test3.txt": {}},
			},
			expected: plan.Plan{
				"test2.txt",
				"test.txt",
				"test3.txt",
			},
		},
		"out of date": {
			folder: "testdata/out_of_date",
			existing: []map[filename]fileContents{
				{"test2.txt": {}},
				{"test.txt": {}},
				{"test3.txt": {}},
				{"test2.txt": {}},
			},
			expected: plan.Plan{
				"test.txt",
				"test3.txt",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			popd := pushd(t, test.folder)
			defer popd()

			for _, batch := range test.existing {
				for f, d := range batch {
					createFile(t, f, d)
				}
			}

			actual := safeDo(t, test.options)

			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Fatalf("plan mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}
