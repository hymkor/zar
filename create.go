package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func addAFile(zw *zip.Writer, thePath string, log io.Writer) ([]string, error) {
	srcFile, err := os.Open(thePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if srcFile != nil {
			srcFile.Close()
		}
	}()
	slashPath := filepath.ToSlash(thePath)

	stat, err := srcFile.Stat()
	if err != nil {
		return nil, err
	}

	var storedFiles []string
	if thePath != "." || thePath != ".." {
		fullpath, err := filepath.Abs(thePath)
		if err != nil {
			return nil, err
		}
		storedFiles = []string{fullpath}
	} else {
		storedFiles = []string{}
	}

	if stat.IsDir() {
		subDir, err := srcFile.Readdir(-1)

		srcFile.Close()
		srcFile = nil

		if err != nil && err != io.EOF {
			return nil, err
		}
		if slashPath[len(slashPath)-1] != '/' {
			slashPath = slashPath + "/"
		}
		_, err = zw.CreateHeader(
			&zip.FileHeader{
				Name:     slashPath,
				NonUTF8:  false,
				Modified: stat.ModTime(),
			})
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(log, "a", slashPath)

		for _, fileInSubDir := range subDir {
			thePath := filepath.Join(thePath, fileInSubDir.Name())
			_storedFiles, err := addAFile(zw, thePath, log)
			if err != nil {
				return nil, err
			}
			storedFiles = append(storedFiles, _storedFiles...)
		}
		return storedFiles, nil
	}

	fileInZipWriter, err := zw.CreateHeader(
		&zip.FileHeader{
			Name:     slashPath,
			NonUTF8:  false,
			Modified: stat.ModTime(),
		})
	if err != nil {
		return nil, err
	}
	io.Copy(fileInZipWriter, srcFile)
	fmt.Fprintln(log, "a", slashPath)
	return storedFiles, nil
}

func create(zipName string, files []string, verbose bool, log io.Writer) ([]string, error) {
	if !verbose {
		log = io.Discard
	}
	succeeded := false
	var w io.Writer
	if zipName == "-" {
		w = os.Stdout
	} else {
		_zipName := zipName + ".tmp"
		_w, err := os.Create(_zipName)
		if err != nil {
			return nil, err
		}
		defer func() {
			_w.Close()
			if succeeded {
				os.Rename(_zipName, zipName)
			}
		}()
		w = _w
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	storedFiles := make([]string, 0)
	for len(files) > 0 {
		if len(files) >= 2 && files[0] == "-C" {
			// -C dir
			if err := os.Chdir(files[1]); err != nil {
				return nil, err
			}
			files = files[2:]
		} else if strings.HasPrefix(files[0], "-C") {
			// -Cdir
			if err := os.Chdir(files[0][2:]); err != nil {
				return nil, err
			}
			files = files[1:]
		} else {
			_storedFiles, err := addAFile(zw, files[0], log)
			if err != nil {
				return nil, err
			}
			files = files[1:]
			storedFiles = append(storedFiles, _storedFiles...)
		}
	}
	succeeded = true
	return storedFiles, nil
}
