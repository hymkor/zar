package main

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/getwild"
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

var (
	flag        = NewFlagSet()
	flagCreate  = flag.Bool("c", false, "Create")
	flagTest    = flag.Bool("t", false, "Test")
	flagExtract = flag.Bool("x", false, "Extract")
	flagVerbose = flag.Bool("v", false, "Verbose")
	flagFile    = flag.String("f", "-", "Filename")
	flagMove    = flag.Bool("remove-files", false, "RemoveFiles")
	flagMd5     = flag.Bool("md5", false, "Show MD5SUM")
)

func mains() error {
	flag.Ignore("C")

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}

	if !isChoosedOne(*flagTest, *flagExtract, *flagCreate) {
		return errors.New("zar.exe: Must specify one of -c, -t, -x")
	}

	if *flagTest {
		return list(*flagFile, flag.Args(), *flagVerbose, os.Stdout)
	} else if *flagExtract {
		return extract(*flagFile, flag.Args(), *flagVerbose, os.Stderr)
	} else if *flagCreate {
		storedFiles, err := create(*flagFile, flag.Args(), *flagVerbose, os.Stderr)
		if err == nil && *flagMove {
			curdir := ""
			for i := len(storedFiles) - 1; i >= 0; i-- {
				if storedFiles[i].CurDir != curdir {
					curdir = storedFiles[i].CurDir
					os.Chdir(curdir)
					println("chdir", curdir)
				}
				println("remove", storedFiles[i].Path)
				os.Remove(storedFiles[i].Path)
			}
		}
		return err
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
