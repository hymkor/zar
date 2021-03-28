package main

import (
	"io"
	"os"
)

type TmpFile struct {
	*os.File
	name string
}

func NewTmpFile() (*TmpFile, error) {
	tmpf, err := os.CreateTemp("", "zar")
	if err != nil {
		return nil, err
	}
	return &TmpFile{File: tmpf, name: tmpf.Name()}, nil
}

func (tmpFile *TmpFile) Close() error {
	if err := tmpFile.File.Close(); err != nil {
		return err
	}
	// println("Remove tmpFile:", tmpFile.Name())
	return os.Remove(tmpFile.Name())
}

func stdin2tmpfile() (*TmpFile, int64, error) {
	tmpf, err := NewTmpFile()
	if err != nil {
		return nil, 0, err
	}
	size, err := io.Copy(tmpf, os.Stdin)
	if err != nil {
		tmpf.Close()
		return nil, 0, err
	}
	if _, err := tmpf.Seek(0, os.SEEK_SET); err != nil {
		tmpf.Close()
		return nil, 0, err
	}
	return tmpf, size, nil
}
