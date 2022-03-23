package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func hasDriveLetter(path string) bool {
	if len(path) < 2 {
		return false
	}
	if path[1] != ':' {
		return false
	}
	if 'A' <= path[0] && path[0] <= 'Z' {
		return true
	}
	if 'a' <= path[0] && path[0] <= 'z' {
		return true
	}
	return false
}

func stripDriveLetterAndRoot(path string) string {
	if hasDriveLetter(path) {
		path = path[2:]
	}
	if len(path) >= 1 && path[0] == '/' {
		path = path[1:]
	}
	return path
}

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
	slashPath := stripDriveLetterAndRoot(filepath.ToSlash(thePath))

	stat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if thePath != "." && thePath != ".." {
		fullpath, err := filepath.Abs(thePath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			fullpath += string(os.PathSeparator)
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

func yesNo(fname string) (string, error) {
	var answer string
	fmt.Printf("overwrite %s? [y]es, [n]o: ", fname)
	_, err := fmt.Scanln(&answer)
	return answer, err
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
		if _, err := os.Stat(zipName); err == nil {
			answer, err := yesNo(zipName)
			if err != nil {
				return err
			}
			if answer != "y" && answer != "Y" {
				return fmt.Errorf("%s: canceled overwriting", zipName)
			}
			os.Rename(zipName, zipName+"~")
		}
		_w, err := os.CreateTemp(filepath.Dir(zipName), "zar*")
		if err != nil {
			return err
		}
		defer func() {
			_w.Close()
			// println("remove", _w.Name())
			if succeeded {
				os.Rename(_w.Name(), zipName)
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
			if strings.IndexAny(files[0], "*?") >= 0 {
				expandedFiles, err := filepath.Glob(files[0])
				if err != nil {
					return fmt.Errorf("%s: Glob: %w", files[0], err)
				}
				for _, fn := range expandedFiles {
					err := addAFile(zw, fn, log, pushStoredFile)
					if err != nil {
						return fmt.Errorf("%s: addAFile: %w", fn, err)
					}
				}
			} else {
				err := addAFile(zw, files[0], log, pushStoredFile)
				if err != nil {
					return fmt.Errorf("%s: AddAFile: %w", files[0], err)
				}
			}
			files = files[1:]
		}
	}
	succeeded = true
	return nil
}
