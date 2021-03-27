package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func List(fileName string, verbose bool, w io.Writer) error {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if verbose {
			mode := f.Mode()
			bit := fs.FileMode(01000)
			for _, c := range "drwxrwxrwx" {
				if (mode & bit) > 0 {
					fmt.Fprintf(w, "%c", c)
				} else {
					fmt.Fprint(w, "-")
				}
				bit >>= 1
			}
			fmt.Fprintf(w, "%8d %s ",
				f.CompressedSize64,
				f.Modified.Format("01 _2 15:04"))
		}
		fmt.Fprintln(w, f.Name)
	}
	return nil
}

func mains() error {
	flag := NewFlagSet()
	flagTest := flag.Bool("t", false, "Test")
	flagFile := flag.String("f", "-", "Filename")
	flagVerbose := flag.Bool("v", false, "Verbose")

	if err := flag.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *flagTest {
		return List(*flagFile, *flagVerbose, os.Stdout)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
