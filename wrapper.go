package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/zetamatta/go-windows-mbcs"
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

func doEach(fileName string, files []string, f func(name string, sc *ZipScanner) error) error {
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
		if err := f(name, sc); err != nil {
			return err
		}
	}
	return nil
}
