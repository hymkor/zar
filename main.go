package main

import (
	"errors"
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
	var (
		flag        = NewFlagSet()
		flagCreate  = flag.Bool("c", false, "Create")
		flagTest    = flag.Bool("t", false, "Test")
		flagExtract = flag.Bool("x", false, "Extract")
		flagVerbose = flag.Bool("v", false, "Verbose")
		flagFile    = flag.String("f", "-", "Filename")
	)

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}

	if !isChoosedOne(*flagTest, *flagExtract, *flagCreate) {
		return errors.New("Choose one of -c,-t and -x")
	}

	if *flagTest {
		return list(*flagFile, flag.Args(), *flagVerbose, os.Stdout)
	} else if *flagExtract {
		return extract(*flagFile, flag.Args(), *flagVerbose, os.Stderr)
	} else if *flagCreate {
		return create(*flagFile, flag.Args(), *flagVerbose, os.Stderr)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
