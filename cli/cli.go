// Package cli provides functionalities to convert images from some format to another.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeConvertError   = 1
)

type App struct {
	OutStream, ErrStream io.Writer
}

// Run is the entry point to the cli. Parses the arguments slice and routes to the proper flag/args combination
func (a *App) Run(args []string) int {
	inExt, outExt, rootPath, errCode := parseFlags(a.ErrStream, args)
	if errCode != 0 {
		return errCode
	}
	inExt = "." + inExt
	outExt = "." + outExt
	fmt.Println("in: %s, out: %s", inExt, outExt)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if !isValidExtension(ext) {
			return errors.New(fmt.Sprintf("%s is not a valid file", path))
		}
		if ext != inExt {
			return nil
		}
		baseLen := len(path) - len(ext)
		baseName := path[:baseLen]
		outFileName := getNewFileName(baseName, outExt)
		err = convert(path, outFileName)
		return err
	})
	if err != nil {
		fmt.Fprintf(a.ErrStream, "error: %v\n", err)
		return ExitCodeConvertError
	}

	return ExitCodeOK
}

func parseFlags(errStream io.Writer, args []string) (inExt, outExt, rootPath string, errCode int) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(errStream)
	flags.Usage = func() {
		fmt.Fprintf(errStream, "Usage: %s image_dir [options]\n", args[0])
		fmt.Fprintln(errStream, "options:")
		flags.PrintDefaults()
	}

	flags.StringVar(&inExt, "i", "jpg", "input file extension <png, jpeg, jpg>")
	flags.StringVar(&outExt, "o", "png", "output file extension <png, jpeg, jpg>")
	if err := flags.Parse(args[1:]); err != nil {
		errCode = ExitCodeParseFlagError
		return
	}
	rootPath = flags.Arg(0)
	if rootPath == "" {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		errCode = ExitCodeParseFlagError
		return
	}
	if _, err := os.Stat(rootPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "error: %s: no such file or directory\n", rootPath)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		errCode = ExitCodeParseFlagError
		return
	}
	return
}

func isValidExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

func getNewFileName(baseName, outExt string) string {
	filename := baseName + outExt
	for n := 2; ; n++ {
		if _, err := os.Stat(filename); err == nil {
			// Already exists
			filename = fmt.Sprintf("%s (%d)%s", baseName, n, outExt)
			continue
		} else if errors.Is(err, os.ErrNotExist) {
			// Does not exists (new file name)
			break
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			panic("File may or may not exist. os.Stat error")
		}
	}
	return filename
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
		panic(fmt.Sprintf("Unavailable file extension: %s", inFileName))
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
		panic(fmt.Sprintf("Unavailable file extension: %s", outFileName))
	}
	return err
}
