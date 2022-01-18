package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func addAFile(zw *zip.Writer, thePath string, log io.Writer, pushStoredFile func(string)) error {
	srcFile, err := os.Open(thePath)
	if err != nil {
		return err
	}
	defer func() {
		if srcFile != nil {
			srcFile.Close()
		}
	}()
	slashPath := filepath.ToSlash(thePath)

	stat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if thePath != "." || thePath != ".." {
		fullpath, err := filepath.Abs(thePath)
		if err != nil {
			return err
		}
		pushStoredFile(fullpath)
	}

	if stat.IsDir() {
		subDir, err := srcFile.Readdir(-1)

		srcFile.Close()
		srcFile = nil

		if err != nil && err != io.EOF {
			return err
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
			return err
		}
		fmt.Fprintln(log, "a", slashPath)

		for _, fileInSubDir := range subDir {
			thePath := filepath.Join(thePath, fileInSubDir.Name())
			err := addAFile(zw, thePath, log, pushStoredFile)
			if err != nil {
				return err
			}
		}
		return nil
	}

	fileInZipWriter, err := zw.CreateHeader(
		&zip.FileHeader{
			Name:     slashPath,
			NonUTF8:  false,
			Modified: stat.ModTime(),
		})
	if err != nil {
		return err
	}
	io.Copy(fileInZipWriter, srcFile)
	fmt.Fprintln(log, "a", slashPath)
	return nil
}

func create(zipName string, files []string, verbose bool, log io.Writer, pushStoredFile func(string)) error {
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
			return err
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
			err := addAFile(zw, files[0], log, pushStoredFile)
			if err != nil {
				return err
			}
			files = files[1:]
		}
	}
	succeeded = true
	return nil
}
