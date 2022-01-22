package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/getwild"

	"github.com/zetamatta/zar/internal/stringstack"
)

var version string

func programInfo() string {
	name, err := os.Executable()
	if err != nil {
		name = "zar"
	}
	return fmt.Sprintf("%s %s", name, version)
}

var (
	flagCreate  = false
	flagTest    = false
	flagExtract = false
	flagVerbose = false
	flagFile    = "-"
	flagMove    = false
)

var flagHash func() hash.Hash = nil

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
		return fmt.Errorf("%s\nMust specify one of -c, -t, -x", programInfo())
	}
	parsedOptionCount := 0
	for len(args) > 0 && len(args[0]) > 0 && args[0][0] == '-' {
		parsedOptionCount++
		if len(args[0]) >= 2 && args[0][1] == '-' {
			switch args[0] {
			case "--md5":
				flagHash = md5.New
			case "--sha1":
				flagHash = sha1.New
			case "--sha256":
				flagHash = sha256.New
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
	if parsedOptionCount <= 0 {
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
		var fnameStack stringstack.Stack
		push := func(string) {}

		if flagMove {
			push = func(fn string) {
				fnameStack.Push(fn)
			}
		}
		err := create(flagFile, args, flagVerbose, os.Stderr, push)

		if err == nil && flagMove {
			var buffer strings.Builder
			for fnameStack.PopTo(&buffer) {
				thePath := buffer.String()

				switch thePath[len(thePath)-1] {
				case '/', '\\':
					fmt.Fprintln(os.Stderr, "rmdir", thePath)
				default:
					fmt.Fprintln(os.Stderr, "rm", thePath)
				}
				if thePath == "." || thePath == ".." {
					continue
				}
				os.Remove(thePath)

				buffer.Reset()
			}
		}
		return err
	} else {
		return fmt.Errorf("%s\nMust specify one of -c, -t, -x", programInfo())
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
