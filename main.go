package main

import (
	"fmt"
	"os"
)

func mains() error {
	flag := NewFlagSet()
	flagTest := flag.Bool("t", false, "Test")
	flagFile := flag.String("f", "-", "Filename")
	flagVerbose := flag.Bool("v", false, "Verbose")

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *flagTest {
		return List(*flagFile, *flagVerbose, os.Stdout)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
