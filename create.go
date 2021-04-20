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
			addAFile(zw, thePath, log)
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
	w, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer w.Close()

	zw := zip.NewWriter(w)
	defer zw.Close()

	for _, name := range files {
		addAFile(zw, name, log)
	}
	return nil
}
