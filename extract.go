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
			fmt.Fprintln(log, dir)
			w, err = os.Create(name)
			if err != nil {
				return err
			}
		}
		fmt.Fprintln(log, name)
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
