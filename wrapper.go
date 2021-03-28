package main

import (
	"archive/zip"
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
