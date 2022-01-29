package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func yesNoAllNoneQuit(fname string) (string, error) {
	var answer string
	// fmt.Printf("replace %s? [y]es, [n]o, [A]ll, [N]one, [r]ename: ", fname)
	fmt.Printf("replace %s? [y]es, [n]o, [A]ll, [N]one: ", fname)
	_, err := fmt.Scanln(&answer)
	return answer, err
}

func extract(fileName string, files []string, verbose bool, log io.Writer) error {
	if !verbose {
		log = io.Discard
	}
	openModeDefault := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	alwaysNone := false
	return doEach(fileName, files, func(name string, sc *ZipScanner) error {
		name = filepath.FromSlash(stripDriveLetterAndRoot(name))
		fileInfo := sc.FileInfo()
		if fileInfo.IsDir() {
			fmt.Fprintln(log, "mkdir", name)
			mode := fileInfo.Mode()
			return os.MkdirAll(name, mode)
		}
		w, err := os.OpenFile(name, openModeDefault, sc.Mode())
		if err != nil {
			if os.IsExist(err) {
				if alwaysNone {
					return nil
				}
				for {
					answer, err := yesNoAllNoneQuit(name)
					if err != nil {
						return err
					}
					if answer == "y" {
						break
					} else if answer == "a" || answer == "A" {
						openModeDefault &^= os.O_EXCL
						break
					} else if answer == "n" {
						return nil
					} else if answer == "N" {
						alwaysNone = true
						return nil
					}
				}
			} else {
				dir := filepath.Dir(name)
				if err := os.MkdirAll(dir, 0777); err != nil {
					return err
				}
				fmt.Fprintln(log, "mkdir", dir)
			}
			w, err = os.OpenFile(name, openModeDefault&^os.O_EXCL, sc.Mode())
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
