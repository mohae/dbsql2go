package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSetOutput(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "dbsql2go")
	if err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	gopath := os.Getenv("$GOPATH")
	if gopath == "" {
		gopath = "$HOME/go"
	}
	gopath = filepath.Join(os.ExpandEnv(gopath), "src")
	dbName = "dbsql2go_test"
	outputfile := dbName + ".go"

	tests := []struct {
		out          string
		filePerTable bool
		expOut       string
		expFilename  string
	}{
		{"", false, wd, outputfile},
		{"stdout", false, "stdout", ""},
		{"out.go", false, ".", "out.go"},
		{filepath.Join(tmpDir, "test"), false, filepath.Join(tmpDir, "test"), outputfile},
		{filepath.Join(tmpDir, "out.go"), false, tmpDir, "out.go"},
		{filepath.Join(tmpDir, "test", "out.go"), false, filepath.Join(tmpDir, "test"), "out.go"},
		{gopath, false, gopath, "dbsql2go_test.go"},
		{filepath.Join(gopath, "test"), false, filepath.Join(gopath, "test"), "dbsql2go_test.go"},
		{filepath.Join(gopath, "out.go"), false, gopath, "out.go"},
		{filepath.Join(gopath, "test", "out.go"), false, filepath.Join(gopath, "test"), "out.go"},

		{"", true, wd, outputfile},
		{"stdout", true, "stdout", ""},
		{"out.go", true, ".", "out.go"},
		{filepath.Join(tmpDir, "test"), true, filepath.Join(tmpDir, "test"), outputfile},
		{filepath.Join(tmpDir, "out.go"), true, tmpDir, "out.go"},
		{filepath.Join(tmpDir, "test", "out.go"), true, filepath.Join(tmpDir, "test"), "out.go"},
		{gopath, true, gopath, "dbsql2go_test.go"},
		{filepath.Join(gopath, "test"), true, filepath.Join(gopath, "test"), "dbsql2go_test.go"},
		{filepath.Join(gopath, "out.go"), true, gopath, "out.go"},
		{filepath.Join(gopath, "test", "out.go"), true, filepath.Join(gopath, "test"), "out.go"},
	}

	for i, test := range tests {
		out = test.out
		_, fname, err := setOutput()
		if err != nil {
			t.Errorf("%d:%s: unexpected error: %q", i, test.out, err)
			continue
		}
		if out != test.expOut {
			t.Errorf("%d:%s: got %q want %q", i, test.out, out, test.expOut)
		}
		if fname != test.expFilename {
			t.Errorf("%d:%s: got %q want %q", i, test.out, fname, test.expFilename)
		}
		if test.out != "stdout" {
			os.Remove(filepath.Join(test.expOut, fname))
		}
	}
}
