package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func makeMatchingFunc(files []string) func(string) bool {
	if files == nil || len(files) <= 0 {
		return func(string) bool { return true }
	}
	_files := make([]string, len(files))
	for i, f := range files {
		_files[i] = filepath.ToSlash(f)
	}
	sort.Strings(_files)

	return func(name string) bool {
		index := sort.Search(len(_files), func(i int) bool {
			if m, err := path.Match(_files[i], name); err == nil && m {
				return true
			}
			return _files[i] >= name
		})
		if index < 0 || index >= len(_files) {
			return false
		}
		m, err := path.Match(_files[index], name)
		return err == nil && m
	}
}

func List(fileName string, files []string, verbose bool, w io.Writer) error {
	zipReader, err := NewZipReadWrapper(fileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	isMatch := makeMatchingFunc(files)

	sc := zipReader.NewScanner()
	for sc.Scan() {
		name, err := sc.Name()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if !isMatch(name) {
			continue
		}
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
		fmt.Fprintln(w, name)
	}
	return nil
}
