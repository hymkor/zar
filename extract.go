package main

import (
	"fmt"
	"io"
	"os"

	"path/filepath"
)

func extract(fileName string, files []string, openMode int, verbose bool, log io.Writer) error {
	if !verbose {
		log = io.Discard
	}
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		name = filepath.FromSlash(stripDriveLetterAndRoot(name))
		fileInfo := sc.FileInfo()
		if fileInfo.IsDir() {
			fmt.Fprintln(log, "mkdir", name)
			mode := fileInfo.Mode()
			return os.MkdirAll(name, mode)
		}
		w, err := os.OpenFile(name, openMode, sc.Mode())
		if err != nil {
			if os.IsExist(err) {
				return err
			}
			dir := filepath.Dir(name)
			if err := os.MkdirAll(dir, 0777); err != nil {
				return err
			}
			fmt.Fprintln(log, "mkdir", dir)
			w, err = os.OpenFile(name, openMode, sc.Mode())
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
