package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func addAFile(zw *zip.Writer, name string, log io.Writer) error {
	fmt.Fprintln(log, name)

	srcFile, err := os.Open(name)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	stat, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		subDir, err := srcFile.Readdir(-1)
		if err != nil && err != io.EOF {
			return err
		}
		for _, fileInSubDir := range subDir {
			thePath := filepath.Join(name, fileInSubDir.Name())
			if err := addAFile(zw, thePath, log); err != nil {
				return err
			}
		}
		return nil
	}
	fileInZipWriter, err := zw.CreateHeader(
		&zip.FileHeader{
			Name:     filepath.ToSlash(name),
			NonUTF8:  false,
			Modified: stat.ModTime(),
		})
	if err != nil {
		return err
	}
	io.Copy(fileInZipWriter, srcFile)
	return nil
}

func addAfileOn(zw *zip.Writer, name string, log io.Writer, dir string) error {
	origDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}
	defer os.Chdir(origDir)
	return addAFile(zw, name, log)
}

func create(zipName string, files []string, verbose bool, log io.Writer) error {
	if !verbose {
		log = io.Discard
	}
	var w io.Writer
	if zipName == "-" {
		w = os.Stdout
	} else {
		_w, err := os.Create(zipName)
		if err != nil {
			return err
		}
		defer _w.Close()
		w = _w
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	for len(files) > 0 {
		var err error
		if len(files) >= 3 && files[0] == "-C" {
			// -C dir file
			err = addAfileOn(zw, files[2], log, files[1])
			files = files[3:]
		} else if len(files) >= 2 && strings.HasPrefix(files[0], "-C") {
			// -Cdir file
			err = addAfileOn(zw, files[1], log, files[0][2:])
			files = files[2:]
		} else {
			err = addAFile(zw, files[0], log)
			files = files[1:]
		}
		if err != nil {
			return err
		}
	}
	return nil
}
