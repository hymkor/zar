package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nyaosorg/go-windows-mbcs"
)

type ZipReadWrapper struct {
	*zip.Reader
	closer func() error
}

func (z *ZipReadWrapper) Close() error {
	// println("ZipReadWrapper.Close()")
	return z.closer()
}

func NewZipReadWrapper(fileName string) (*ZipReadWrapper, error) {
	if fileName == "-" {
		tmpf, size, err := stdin2tmpfile()
		if err != nil {
			return nil, err
		}
		zipReader, err := zip.NewReader(tmpf, size)
		if err != nil {
			tmpf.Close()
			return nil, err
		}
		return &ZipReadWrapper{
			Reader: zipReader,
			closer: func() error { return tmpf.Close() },
		}, nil
	} else {
		zipReadCloser, err := zip.OpenReader(fileName)
		if err != nil {
			return nil, err
		}
		return &ZipReadWrapper{
			Reader: &zipReadCloser.Reader,
			closer: func() error { return zipReadCloser.Close() },
		}, nil
	}
}

func (z *ZipReadWrapper) NewScanner() *ZipScanner {
	return &ZipScanner{
		zrw:   z,
		File:  nil,
		index: -1,
	}
}

type ZipScanner struct {
	*zip.File
	zrw   *ZipReadWrapper
	index int
}

func (e *ZipScanner) Scan() bool {
	e.index++
	if e.index >= len(e.zrw.File) {
		return false
	}
	e.File = e.zrw.File[e.index]
	return true
}

func (e *ZipScanner) Name() (string, error) {
	f := e.File
	if f.NonUTF8 {
		return mbcs.AtoU([]byte(f.Name), mbcs.ACP)
	} else {
		return f.Name, nil
	}
}

func makeMatchingFunc(pattern []string) func(string) (bool, string) {
	if pattern == nil || len(pattern) <= 0 {
		return func(string) (bool, string) { return true, "" }
	}
	_patterns := make([]string, 0, len(pattern))
	_dirs := make([]string, 0)
	original := make(map[string]string)
	for _, f := range pattern {
		newName := filepath.ToSlash(f)
		if strings.HasSuffix(newName, "/") {
			_dirs = append(_dirs, newName)
		} else {
			_patterns = append(_patterns, newName)
		}
		original[newName] = f
	}
	sort.Strings(_dirs)
	sort.Strings(_patterns)

	return func(name string) (bool, string) {
		// search dirs
		index := sort.Search(len(_dirs), func(i int) bool {
			if strings.HasPrefix(name, _dirs[i]) {
				return true
			}
			return _dirs[i] >= name
		})
		if index >= 0 && index < len(_dirs) {
			if strings.HasPrefix(name, _dirs[index]) {
				return true, original[_dirs[index]]
			}
		}
		// search wildcards
		index = sort.Search(len(_patterns), func(i int) bool {
			if m, err := path.Match(_patterns[i], name); err == nil && m {
				return true
			}
			return _patterns[i] >= name
		})
		if index >= 0 && index < len(_patterns) {
			m, err := path.Match(_patterns[index], name)
			return (err == nil && m), original[_patterns[index]]
		}
		return false, ""
	}
}

func FilterOptionC(args []string) ([]string, map[string]string) {
	files := make([]string, 0, len(args))
	fileToPut := make(map[string]string, len(args))
	dir := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "-C" && i+1 < len(args) {
			dir = args[i+1]
			i++
		} else if strings.HasPrefix(args[i], "-C") {
			dir = args[i][2:]
		} else {
			if dir != "" {
				fileToPut[args[i]] = dir
			}
			files = append(files, args[i])
		}
	}
	return files, fileToPut
}

func doEach(fileName string, files []string, f func(name string, sc *ZipScanner) error) error {
	zipReader, err := NewZipReadWrapper(fileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	var fileToPut map[string]string
	files, fileToPut = FilterOptionC(files)
	isMatch := makeMatchingFunc(files)

	sc := zipReader.NewScanner()
	for sc.Scan() {
		name, err := sc.Name()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		ok, matchedPattern := isMatch(name)
		if !ok {
			continue
		}
		dir := fileToPut[matchedPattern]
		if dir != "" {
			saveDir, _err := os.Getwd()
			if _err != nil {
				return _err
			}
			// println("Chdir:",dir)
			if _err := os.Chdir(dir); _err != nil {
				return _err
			}
			err = f(name, sc)
			if _err := os.Chdir(saveDir); _err != nil {
				return _err
			}
			// println("Chdir:",saveDir)
		} else {
			err = f(name, sc)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
