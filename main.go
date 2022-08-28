package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	i := flag.String("i", "jpg", "output format [png, jpeg, jpg] is available.")
	o := flag.String("o", "png", "output format [png, jpeg, jpg] is available.")
	flag.Parse()
	rootPath := flag.Arg(0)
	if rootPath == "" {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		return
	}
	if _, err := os.Stat(rootPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return
	}
	inExt := "." + *i
	outExt := "." + *o
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext != inExt {
				return errors.New(fmt.Sprintf("%s is not a valid file", path))
			}
			baseLen := len(path) - len(inExt)
			baseName := path[:baseLen]
			outFileName := getNewFileName(baseName, outExt)
			err = convert(path, outFileName)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}

func existsPath(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		// path exists
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// path does not exist
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false, err
	}
}

func getNewFileName(baseName, outExt string) string {
	outFileName := baseName + outExt
	exists, err := existsPath(outFileName)
	n := 1
	for err == nil && exists {
		fmt.Println("update")
		n++
		outFileName = fmt.Sprintf("%s (%d)%s", baseName, n, outExt)
		exists, err = existsPath(outFileName)
	}
	return outFileName
}

func convert(inFileName, outFileName string) error {
	// open input file
	fin, err := os.Open(inFileName)
	if err != nil {
		return err
	}
	defer fin.Close()

	// create output file
	fout, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer fout.Close()

	// Decode
	var img image.Image
	switch filepath.Ext(inFileName) {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(fin)
	case ".png":
		img, err = png.Decode(fin)
	default:
		return errors.New("Unknown extension")
	}
	if err != nil {
		return err
	}

	// Encode
	switch filepath.Ext(outFileName) {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(fout, img, nil)
	case ".png":
		err = png.Encode(fout, img)
	default:
		return errors.New("Unknown extension")
	}
	return err
}
