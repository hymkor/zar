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
		if len(files) >= 2 && files[0] == "-C" {
			// -C dir
			if err := os.Chdir(files[1]); err != nil {
				return err
			}
			files = files[2:]
		} else if strings.HasPrefix(files[0], "-C") {
			// -Cdir
			if err := os.Chdir(files[0][2:]); err != nil {
				return err
			}
			files = files[1:]
		} else {
			if err := addAFile(zw, files[0], log); err != nil {
				return err
			}
			files = files[1:]
		}
	}
	return nil
}
