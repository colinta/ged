package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// testCase represents a single YAML test case.
type testCase struct {
	Name   string   `yaml:"name"`
	Input  string   `yaml:"input"`
	Args   []string `yaml:"args"`
	Output string   `yaml:"output"`
	Error  bool     `yaml:"error"`
}

func TestYAML(t *testing.T) {
	files, err := filepath.Glob("testdata/*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test files found in testdata/")
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("reading %s: %v", file, err)
		}

		var cases []testCase
		if err := yaml.Unmarshal(data, &cases); err != nil {
			t.Fatalf("parsing %s: %v", file, err)
		}

		suite := strings.TrimSuffix(filepath.Base(file), ".yaml")
		t.Run(suite, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					in := strings.NewReader(tc.Input)
					out := &bytes.Buffer{}

					err := run(tc.Args, in, out, io.Discard)

					if tc.Error {
						if err == nil {
							t.Fatal("expected error, got nil")
						}
						return
					}

					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if out.String() != tc.Output {
						t.Errorf("output mismatch\ngot:\n%s\nwant:\n%s", out.String(), tc.Output)
					}
				})
			}
		})
	}
}
