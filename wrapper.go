package main

import (
	"archive/zip"

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
