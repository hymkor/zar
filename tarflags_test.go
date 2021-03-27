package main

import (
	"testing"
)

func TestFlagSet(t *testing.T) {
	f := NewFlagSet()
	aFlag := f.Bool("a", false, "usage")
	bFlag := f.String("b", "", "usage")
	cFlag := f.Bool("c", false, "usage")
	vFlag := f.Bool("v", false, "usage")
	fFlag := f.String("f", "", "usage")
	removeFlag := f.Bool("remove-files", false, "usage")
	notSetBool := f.Bool("not-set-bool", false, "usage")
	notSetString := f.String("not-set-string", "", "usage")

	const filename = "foo.tar.gz"
	const PARAM1 = "PARAMETER1"
	const PARAM2 = "PARAMETER2"

	err := f.Parse([]string{"-a", "-b", "BOPTION", "-cvf", "--remove-files", filename, PARAM1, PARAM2})
	if err != nil {
		t.Fatalf("Parse() fails: %s", err.Error())
	}
	if !*aFlag {
		t.Fatal("-a: fails: expect true, but false")
	}
	if *bFlag != "BOPTION" {
		t.Fatalf("-b: fails: expect BOPTION, but %s", *bFlag)
	}
	if !*cFlag {
		t.Fatal("-c: fails: expect true, but false")
	}
	if !*vFlag {
		t.Fatal("-v: fails: expect true, but false")
	}
	if *fFlag != filename {
		t.Fatalf("-f: fails: expect \"%s\", but \"%s\"", filename, *fFlag)
	}
	if !*removeFlag {
		t.Fatal("--remove-files: fails: expect true, but false")
	}
	if *notSetBool {
		t.Fatal("--not-set-bool: fails: expect false, but true")
	}
	if *notSetString != "" {
		t.Fatalf("--not-set-string: fails: expect \"\", but \"%s\"", *notSetString)
	}
	if n := f.NArg(); n != 2 {
		t.Fatalf("f.NArgs(): fails: expect 2, but %d", n)
	}
	if s := f.Arg(0); s != PARAM1 {
		t.Fatalf("f.Arg(0): fails: expect %s, but %s", PARAM1, s)
	}
	if s := f.Arg(1); s != PARAM2 {
		t.Fatalf("f.Arg(0): fails: expect %s, but %s", PARAM2, s)
	}
}
