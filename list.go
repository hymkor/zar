package main

import (
	"fmt"
	"io"
	"io/fs"
)

func list(fileName string, files []string, verbose bool, w io.Writer) error {
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		if verbose {
			if sc.FileInfo().IsDir() {
				w.Write([]byte{'d'})
			} else {
				w.Write([]byte{'-'})
			}
			mode := sc.Mode()
			bit := fs.FileMode(01000)
			for _, c := range []byte("rwxrwxrwx") {
				bit >>= 1
				if (mode & bit) > 0 {
					w.Write([]byte{c})
				} else {
					w.Write([]byte{'-'})
				}
			}
			fmt.Fprintf(w, "%8d %s ",
				sc.CompressedSize64,
				sc.Modified.Format("2006/01/02 15:04"))
			if sc.NonUTF8 {
				fmt.Fprint(w, "A ")
			} else {
				fmt.Fprint(w, "U ")
			}
		}
		fmt.Fprintln(w, name)
		return nil
	})
}
