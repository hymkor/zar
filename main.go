package main

import (
	"fmt"
	"os"
)

func isChoosedOne(flags ...bool) bool {
	count := 0
	for _, f := range flags {
		if f {
			count++
		}
	}
	return count == 1
}

func mains() error {
	flag := NewFlagSet()
	flagTest := flag.Bool("t", false, "Test")
	flagFile := flag.String("f", "-", "Filename")
	flagVerbose := flag.Bool("v", false, "Verbose")

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *flagTest {
		return List(*flagFile, flag.Args(), *flagVerbose, os.Stdout)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
