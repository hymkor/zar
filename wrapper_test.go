package main_test

import (
	"testing"

	"github.com/zetamatta/zar"
)

func TestFilterOptionC(t *testing.T) {
	source := []string{"-C", "dir", "file1", "file2"}
	file, fileToPut := main.FilterOptionC(source)
	if len(file) != 2 || file[0] != "file1" || file[1] != "file2" {
		t.Fatalf("file failed: %v", file)
	}
	if dir := fileToPut["file1"]; dir != "dir" {
		t.Fatalf("expect '%v' but '%v'", "dir", dir)
	}

	source = []string{"-Cdir", "file1", "file2"}
	file, fileToPut = main.FilterOptionC(source)
	if len(file) != 2 || file[0] != "file1" || file[1] != "file2" {
		t.Fatalf("file failed: %v", file)
	}
	if dir := fileToPut["file1"]; dir != "dir" {
		t.Fatalf("expect '%v' but '%v'", "dir", dir)
	}

}
