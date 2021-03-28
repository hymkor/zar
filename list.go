package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/zetamatta/go-windows-mbcs"
)

func List(fileName string, verbose bool, w io.Writer) error {
	zipReader, err := NewZipReadWrapper(fileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if verbose {
			mode := f.Mode()
			bit := fs.FileMode(01000)
			for _, c := range "drwxrwxrwx" {
				if (mode & bit) > 0 {
					fmt.Fprintf(w, "%c", c)
				} else {
					fmt.Fprint(w, "-")
				}
				bit >>= 1
			}
			fmt.Fprintf(w, "%8d %s ",
				f.CompressedSize64,
				f.Modified.Format("01 _2 15:04"))
		}
		if f.NonUTF8 {
			utf8, err := mbcs.AtoU([]byte(f.Name), mbcs.ACP)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			fmt.Fprintln(w, utf8)
		} else {
			fmt.Fprintln(w, f.Name)
		}
	}
	return nil
}
