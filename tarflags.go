package main

import (
	"fmt"
	"strings"
)

type FlagSet struct {
	bools     map[string]*bool
	strings   map[string]*string
	arguments []string
}

func NewFlagSet() *FlagSet {
	return &FlagSet{
		bools:     map[string]*bool{},
		strings:   map[string]*string{},
		arguments: nil,
	}
}

func (f *FlagSet) Bool(name string, value bool, _ string) *bool {
	f.bools[name] = &value
	return &value
}

func (f *FlagSet) String(name string, value string, usage string) *string {
	f.strings[name] = &value
	return &value
}

func (f *FlagSet) Parse(arguments []string) error {
	options := []string{}
	args := []string{}
	for _, arg1 := range arguments {
		if len(arg1) > 1 && arg1[0] == '-' {
			options = append(options, arg1[1:])
		} else {
			args = append(args, arg1)
		}
	}
	if len(options) == 0 && len(args) > 0 {
		options = append(options, args[0])
		args = args[1:]
	}
	for _, opt := range options {
		if len(opt) > 0 && opt[0] == '-' { // long option
			// fmt.Printf("Long Option: '%s'\n", opt)
			name := opt[1:]
			if store, ok := f.bools[name]; ok {
				*store = true
				continue
			}
			eq := strings.SplitN(name, "=", 2)
			store, ok := f.strings[eq[0]]
			if !ok {
				return fmt.Errorf("'%s': no such options", eq[0])
			}
			if len(eq) >= 2 {
				*store = eq[1]
			} else {
				if len(args) < 1 {
					return fmt.Errorf("'%s': too few arguments", eq[0])
				}
				*store = args[0]
				args = args[1:]
			}
		} else { // short option
			// fmt.Printf("Short Option: '%s'\n", opt)
			for _, c := range opt {
				s := string(c)
				if store, ok := f.bools[s]; ok {
					// fmt.Printf("Found '%s'\n", s)
					*store = true
					continue
				}
				store, ok := f.strings[s]
				if !ok {
					return fmt.Errorf("'%s': no such options", s)
				}
				if len(args) < 1 {
					return fmt.Errorf("'%s': too few arguments", s)
				}
				*store = args[0]
				args = args[1:]
			}
		}
	}
	f.arguments = args
	return nil
}

func (f *FlagSet) Args() []string {
	return f.arguments
}

func (f *FlagSet) Arg(i int) string {
	return f.arguments[i]
}

func (f *FlagSet) NArg() int {
	return len(f.arguments)
}
