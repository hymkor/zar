package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/getwild"
)

var (
	flagCreate  = false
	flagTest    = false
	flagExtract = false
	flagVerbose = false
	flagFile    = "-"
	flagMove    = false
	flagMd5     = false
)

func parseShortOption(flags string, args []string) ([]string, error) {
	for i, c := range flags {
		switch c {
		case 'c':
			flagCreate = true
		case 't':
			flagTest = true
		case 'x':
			flagExtract = true
		case 'v':
			flagVerbose = true
		case 'f':
			var fname string
			if i+1 < len(flags) {
				fname = flags[i+1:]
			} else if len(args) > 0 {
				fname = args[0]
				args = args[1:]
			} else {
				return nil, errors.New("-f without filename")
			}
			var err error
			if flagFile, err = filepath.Abs(fname); err != nil {
				return nil, err
			}
			return args, nil
		case 'C':
			var curdir string
			if i+1 < len(flags) {
				curdir = flags[i+1:]
			} else if len(args) > 0 {
				curdir = args[0]
				args = args[1:]
			} else {
				return nil, errors.New("-C without directory")
			}
			return args, os.Chdir(curdir)
		default:
			return nil, fmt.Errorf("-%c: unknown option", c)
		}
	}
	return args, nil
}

func mains(args []string) error {
	if len(args) <= 0 {
		return errors.New("zar.exe: Must specify one of -c, -t, -x")
	}
	for len(args) > 0 && len(args[0]) > 0 && args[0][0] == '-' {
		if len(args[0]) >= 2 && args[0][1] == '-' {
			switch args[0] {
			case "--md5":
				flagMd5 = true
			case "--remove-files":
				flagMove = true
			default:
				return fmt.Errorf("%s: unknown option", args[0])
			}
			args = args[1:]
		} else {
			var err error
			args, err = parseShortOption(args[0][1:], args[1:])
			if err != nil {
				return err
			}
		}
	}
	if !flagCreate && !flagTest && !flagExtract {
		var err error
		args, err = parseShortOption(args[0], args[1:])
		if err != nil {
			return err
		}
	}

	if flagTest {
		return list(flagFile, args, flagVerbose, os.Stdout)
	} else if flagExtract {
		return extract(flagFile, args, flagVerbose, os.Stderr)
	} else if flagCreate {
		storedFiles, err := create(flagFile, args, flagVerbose, os.Stderr)
		if err == nil && flagMove {
			for i := len(storedFiles) - 1; i >= 0; i-- {
				fmt.Fprintln(os.Stderr, "remove", storedFiles[i])
				// os.Remove(storedFiles[i])
			}
		}
		return err
	} else {
		return errors.New("zar.exe: Must specify one of -c, -t, -x")
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
