package main

import (
	"fmt"
	"io"
	"os"

	"path/filepath"
)

func extract(fileName string, files []string, verbose bool, log io.Writer) error {
	if !verbose {
		log = io.Discard
	}
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		name = filepath.FromSlash(name)
		fileInfo := sc.FileInfo()
		if fileInfo.IsDir() {
			fmt.Fprintln(log, "mkdir", name)
			mode := fileInfo.Mode()
			return os.MkdirAll(name, mode)
		}
		w, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, sc.Mode())
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			dir := filepath.Dir(name)
			if err := os.MkdirAll(dir, 0777); err != nil {
				return err
			}
			w, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE, sc.Mode())
			fmt.Fprintln(log, "mkdir", dir)
			if err != nil {
				return err
			}
		}
		fmt.Fprintln(log, "x", name)
		r, err := sc.Open()
		if err != nil {
			w.Close()
			return err
		}
		_, err = io.Copy(w, r)
		w.Close()
		r.Close()

		os.Chtimes(name, sc.ModTime(), sc.ModTime())

		return err
	})
}
