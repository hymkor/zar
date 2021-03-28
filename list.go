package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

func List(fileName string, verbose bool, w io.Writer) error {
	zipReader, err := NewZipReadWrapper(fileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	sc := zipReader.NewScanner()
	for sc.Scan() {
		if verbose {
			mode := sc.Mode()
			bit := fs.FileMode(01000)
			for _, c := range []byte("drwxrwxrwx") {
				if (mode & bit) > 0 {
					w.Write([]byte{c})
				} else {
					w.Write([]byte{'-'})
				}
				bit >>= 1
			}
			fmt.Fprintf(w, "%8d %s ",
				sc.CompressedSize64,
				sc.Modified.Format("01 _2 15:04"))
		}
		name, err := sc.Name()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		fmt.Fprintln(w, name)
	}
	return nil
}
