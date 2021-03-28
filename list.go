package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func List(fileName string, verbose bool, w io.Writer) error {
	var zipReader *zip.Reader
	if fileName == "-" {
		tmpf, err := os.CreateTemp("", "zar")
		if err != nil {
			return err
		}
		defer func() {
			tmpf.Close()
			os.Remove(tmpf.Name())
		}()
		size, err := io.Copy(tmpf, os.Stdin)
		if err != nil {
			return err
		}
		if _, err := tmpf.Seek(0, os.SEEK_SET); err != nil {
			return err
		}
		zipReader, err = zip.NewReader(tmpf, size)
		if err != nil {
			return err
		}
	} else {
		zipReadCloser, err := zip.OpenReader(fileName)
		if err != nil {
			return err
		}
		defer zipReadCloser.Close()
		zipReader = &zipReadCloser.Reader
	}
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
		fmt.Fprintln(w, f.Name)
	}
	return nil
}
