package main

import (
	"fmt"
	"io"
)

func extract(fileName string, files []string, verbose bool, w io.Writer) error {
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		if verbose {
			fmt.Fprintln(w, name)
		}
		return nil
	})
}
