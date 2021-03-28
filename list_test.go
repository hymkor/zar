package main

import (
	"testing"
)

func TestMakeMatchingFunc(t *testing.T) {
	isMatch := makeMatchingFunc([]string{
		"a*", "f*", "Z*", "foo/*", "bar\\*",
	})

	var filesInZip = []struct {
		file   string
		expect bool
	}{
		{"a", true},
		{"b", false},
		{"ZZZZ", true},
		{"ccc", false},
		{"foo/", true},
		{"foo/foo", true},
		{"bar/", true},
		{"bar/xxx", true},
	}

	for _, f := range filesInZip {
		if act := isMatch(f.file); act != f.expect {
			t.Fatalf(
				"`%s` must be %v, but %v",
				f.file,
				f.expect,
				act)
		}
	}
}
