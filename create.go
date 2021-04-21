package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	for i := 0; i < len(files); i++ {
		name := files[i]
		if name == "-C" && i+2 < len(files) {
			origDir, err := os.Getwd()
			if err != nil {
				return err
			}
			if err := os.Chdir(files[i+1]); err != nil {
				return err
			}
			if err := addAFile(zw, files[i+2], log); err != nil {
				return err
			}
			if err := os.Chdir(origDir); err != nil {
				return err
			}
			i += 2
		} else {
			if err := addAFile(zw, name, log); err != nil {
				return err
			}
		}
	}
	return nil
}
