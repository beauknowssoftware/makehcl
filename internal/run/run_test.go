package run_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/beauknowssoftware/makehcl/internal/run"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

const (
	local = false
)

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

	tempDir = path.Join(tempDir, t.Name())

	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
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

func safeRun(t *testing.T, o run.DoOptions) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("paniced: %v\n%v", r, string(debug.Stack()))
		}
	}()

	if err := run.Do(o); err != nil {
		err = errors.Wrap(err, "failed to run")
		t.Fatal(err)
	}
}

func readFile(t *testing.T, filename string) fileContents {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		err = errors.Wrapf(err, "failed to read file %v", filename)
		t.Fatal(err)
	}

	return d
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

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

func setEnv(t *testing.T, envs map[string]string) func() {
	original := os.Environ()

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			err = errors.Wrapf(err, "failed to set env %v = %v", k, v)
			t.Fatal(err)
		}
	}

	return func() {
		os.Clearenv()

		for _, e := range original {
			p := strings.SplitN(e, "=", 2)
			if err := os.Setenv(p[0], p[1]); err != nil {
				err = errors.Wrapf(err, "failed to reset env %v = %v", p[0], p[1])
				t.Fatal(err)
			}
		}
	}
}

type filename = string
type fileContents = []byte

func TestRun(t *testing.T) {
	tests := map[string]struct {
		folder    string
		existing  []map[filename]fileContents
		want      map[filename]fileContents
		doNotWant []filename
		options   run.DoOptions
		env       map[string]string
	}{
		"command": {
			folder: "testdata/command",
			want: map[filename]fileContents{
				"test.txt": {},
			},
		},
		"path": {
			folder: "testdata/path",
			existing: []map[filename]fileContents{
				{"test/1/original": fileContents{}},
				{"test/2/original": fileContents{}},
			},
			want: map[filename]fileContents{
				"out/1": {},
				"out/2": {},
			},
		},
		"glob": {
			folder: "testdata/glob",
			existing: []map[filename]fileContents{
				{"1.original": fileContents{}},
				{"2.original": fileContents{}},
			},
			want: map[filename]fileContents{
				"1.txt": {},
				"2.txt": {},
			},
		},
		"dynamic command as": {
			folder: "testdata/dynamic_command_as",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic rule as": {
			folder: "testdata/dynamic_rule_as",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic rule in command": {
			folder: "testdata/dynamic_rule_in_command",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic command alias target": {
			folder: "testdata/dynamic_command_alias_target",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic command default": {
			folder: "testdata/dynamic_command_default",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic rule alias target": {
			folder: "testdata/dynamic_rule_alias_target",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic rule default": {
			folder: "testdata/dynamic_rule_default",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic command": {
			folder: "testdata/dynamic_command",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dynamic rule": {
			folder: "testdata/dynamic_rule",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"complex dynamic rule": {
			folder: "testdata/complex_dynamic_rule",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"touch": {
			folder: "testdata/touch",
			want: map[filename]fileContents{
				"test.txt": {},
			},
		},
		"var": {
			folder: "testdata/var",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"multiple vars": {
			folder: "testdata/multiple_vars",
			want: map[filename]fileContents{
				"test1.txt": fileContents("hello1\n"),
				"test2.txt": fileContents("hello2\n"),
			},
		},
		"dependent vars": {
			folder: "testdata/dependent_vars",
			want: map[filename]fileContents{
				"test.txt.1": fileContents("hello1\n"),
				"test.txt.2": fileContents("hello2\n"),
				"test.txt.3": fileContents("hello3\n"),
				"test.txt.4": fileContents("hello4\n"),
			},
		},
		"single default": {
			folder: "testdata/single_default",
			want: map[filename]fileContents{
				"test.txt": {},
			},
		},
		"tee target": {
			folder: "testdata/tee_target",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"echo": {
			folder: "testdata/echo",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"custom shell": {
			folder: "testdata/custom_shell",
			want: map[filename]fileContents{
				"test.txt": {},
			},
		},
		"env var": {
			folder: "testdata/env_var",
			env: map[string]string{
				"TARGET": "test.txt",
			},
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"command env": {
			folder: "testdata/command_env",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"global env": {
			folder: "testdata/global_env",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"rule env": {
			folder: "testdata/rule_env",
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"corelated env var": {
			folder: "testdata/corelated_env_var",
			env: map[string]string{
				"NAME":  "test",
				"EXT":   "txt",
				"VALUE": "hello",
			},
			want: map[filename]fileContents{
				"test_new.txt": fileContents("hello_new\n"),
			},
		},
		"env": {
			folder: "testdata/env",
			env: map[string]string{
				"TARGET": "test.txt",
			},
			want: map[filename]fileContents{
				"test.txt": fileContents("hello\n"),
			},
		},
		"env default goal second": {
			folder: "testdata/env_default_goal",
			want: map[filename]fileContents{
				"test2.txt": {},
			},
			env: map[string]string{
				"SECOND_TARGET": "true",
			},
			doNotWant: []filename{"test.txt"},
		},
		"env default goal first": {
			folder: "testdata/env_default_goal",
			want: map[filename]fileContents{
				"test.txt": {},
			},
			doNotWant: []filename{"test2.txt"},
		},
		"default goal": {
			folder: "testdata/default_goal",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
			doNotWant: []filename{"test3.txt"},
		},
		"command dependencies": {
			folder: "testdata/command_dependencies",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"non target dependency": {
			folder: "testdata/non_target_dependency",
			existing: []map[filename]fileContents{
				{"test2.txt": fileContents{}},
			},
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"dependencies": {
			folder: "testdata/dependencies",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
			},
		},
		"build once": {
			folder: "testdata/build_once",
			want: map[filename]fileContents{
				"test.txt":  {},
				"test2.txt": {},
				"test3.txt": fileContents("hello\n"),
			},
		},
		"build outdated": {
			folder: "testdata/build_outdated",
			existing: []map[filename]fileContents{
				{"test.txt": fileContents("hello\n")},
				{"test2.txt": fileContents("hello\n")},
			},
			want: map[filename]fileContents{
				"test.txt":  fileContents("hello\nhello\n"),
				"test2.txt": fileContents("hello\nhello2\n"),
				"test3.txt": fileContents("hello3\n"),
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

			resetEnv := setEnv(t, test.env)
			defer resetEnv()

			safeRun(t, test.options)

			for filename, expected := range test.want {
				actual := readFile(t, filename)
				if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
					t.Fatalf("%v mismatch (-want,+got):\n%s", filename, diff)
				}
			}

			for _, filename := range test.doNotWant {
				if fileExists(filename) {
					t.Fatalf("expected %v to not exist", filename)
				}
			}
		})
	}
}
