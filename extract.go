package main

import (
	"fmt"
	"io"
	"os"

	"path/filepath"
)

func extract(fileName string, files []string, verbose bool, w io.Writer) error {
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		if verbose {
			fmt.Fprintln(w, name)
		}
		name = filepath.FromSlash(name)
		fileInfo := sc.FileInfo()
		if fileInfo.IsDir() {
			mode := fileInfo.Mode()
			return os.MkdirAll(name, mode)
		}
		w, err := os.Create(name)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			dir := filepath.Dir(name)
			if err := os.MkdirAll(dir, 0777); err != nil {
				return err
			}
			w, err = os.Create(name)
			if err != nil {
				return err
			}
		}
		r, err := sc.Open()
		if err != nil {
			w.Close()
			return err
		}
		_, err = io.Copy(w, r)
		w.Close()
		r.Close()
		return err
	})
}
