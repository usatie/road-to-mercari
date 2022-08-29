// Package cli provides functionalities to convert images from some format to another.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeConvertError   = 1
)

// App consists of output/error streams.
type App struct {
	OutStream, ErrStream io.Writer
}

// Run is the entry point to the cli. Parses the arguments slice and routes to the proper flag/args combination
func (a *App) Run(args []string) int {
	arg, err := parseFlags(a.ErrStream, args)
	if err != nil {
		return ExitCodeParseFlagError
	}
	inExt := "." + arg.inExt
	outExt := "." + arg.outExt
	if arg.verbose {
		fmt.Println(arg)
	}
	var cnt int
	err = filepath.Walk(arg.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		} else if ext := filepath.Ext(path); ext != inExt {
			return errors.New(fmt.Sprintf("%s is not a valid file", path))
		}

		var outPath string
		if arg.dir != "" {
			outPath = getOutputPath(filepath.Join(arg.dir, filepath.Base(path)), inExt, outExt)
		} else {
			outPath = getOutputPath(path, inExt, outExt)
		}
		if arg.verbose {
			fmt.Printf("%s ---> %s\n", path, outPath)
		}
		err = convert(path, outPath, arg.decoder, arg.encoder)
		if err == nil {
			cnt++
		}
		return err
	})
	if arg.verbose {
		fmt.Printf("\n\nconverted %d files\n", cnt)
	}
	if err != nil {
		fmt.Fprintf(a.ErrStream, "error: %v\n", err)
		return ExitCodeConvertError
	}

	return ExitCodeOK
}

// getOutputPath returns the output file name which does not overwrite existing files..
func getOutputPath(path, inExt, outExt string) string {
	baseLen := len(path) - len(inExt)
	baseName := path[:baseLen]
	outPath := baseName + outExt
	for n := 2; ; n++ {
		if _, err := os.Stat(outPath); err == nil {
			// Already exists
			outPath = fmt.Sprintf("%s (%d)%s", baseName, n, outExt)
			continue
		} else if errors.Is(err, os.ErrNotExist) {
			// Does not exists (new file name)
			break
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			panic("File may or may not exist. os.Stat error")
		}
	}
	return outPath
}

// converts input file format to output file format.
func convert(inPath, outPath string, decoder Decoder, encoder Encoder) error {
	// open input file
	fin, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer fin.Close()

	// create output file
	fout, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	// Decode
	img, err := decoder.Decode(fin)
	if err != nil {
		return err
	}

	// Encode
	err = encoder.Encode(fout, img)
	return err
}
