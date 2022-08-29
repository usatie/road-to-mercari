package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// cliArgument is a collection of available options
type cliArgument struct {
	inExt    string
	outExt   string
	rootPath string
	quality  int
	dir      string
	verbose  bool
	decoder  Decoder
	encoder  Encoder
}

func (c cliArgument) String() string {
	return fmt.Sprintf("Options:\n  inExt: %s\n  outExt: %s\n  rootPath: %s\n  quality: %d\n  dir: %s\n  verbose: %v", c.inExt, c.outExt, c.rootPath, c.quality, c.dir, c.verbose)
}

const usageOptions = `options:
  -i, --input	input file extension <png, jpeg, jpg>
  -o, --output	output file extension <png, jpeg, jpg>
  -d, --dir		directory to which output files
  -q, --quality output encoding quality. higher is better. <1-100>
  -v, --verbose verbose output`

// parseFlags parse the arguments and writes errors to errStream
func parseFlags(errStream io.Writer, args []string) (cliArgument, error) {
	var arg cliArgument
	// Custom flag set to set custom output and usage
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(errStream)
	fs.Usage = func() {
		fmt.Fprintf(errStream, "Usage: %s [options] image_dir\n", args[0])
		fmt.Fprintln(errStream, usageOptions)
	}
	// -i, --input
	fs.StringVar(&arg.inExt, "i", "jpg", "input file extension <png, jpeg, jpg>")
	fs.StringVar(&arg.inExt, "input", "jpg", "input file extension <png, jpeg, jpg>")
	// -o, --output
	fs.StringVar(&arg.outExt, "o", "png", "output file extension <png, jpeg, jpg>")
	fs.StringVar(&arg.outExt, "output", "png", "output file extension <png, jpeg, jpg>")
	// -d, --dir
	fs.StringVar(&arg.dir, "d", "", "directory to which output files")
	fs.StringVar(&arg.dir, "dir", "", "directory to which output files")
	// -q, --quality
	fs.IntVar(&arg.quality, "q", 75, "output encoding quality. higher is better. <1-100>")
	fs.IntVar(&arg.quality, "quality", 75, "output encoding quality. higher is better. <1-100>")
	// -v, --verbose
	fs.BoolVar(&arg.verbose, "v", false, "verbose output")
	fs.BoolVar(&arg.verbose, "verbose", false, "verbose output")
	// Parse
	if err := fs.Parse(args[1:]); err != nil {
		return arg, errors.New("Parse Error")
	}
	// rootPath must be passed
	arg.rootPath = fs.Arg(0)
	if arg.rootPath == "" {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		return arg, errors.New("rootPath is empty")
	}
	// rootPath must exist
	if _, err := os.Stat(arg.rootPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "error: %s: no such file or directory\n", arg.rootPath)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		return arg, err
	}
	// Decoder
	var err error
	arg.decoder, err = newDecoder(arg.inExt)
	if err != nil {
		fmt.Fprintf(errStream, "error: %v\n", err)
		return arg, err
	}
	// Encoder
	arg.encoder, err = newEncoder(arg.outExt, arg.quality)
	if err != nil {
		fmt.Fprintf(errStream, "error: %v\n", err)
		return arg, err
	}
	return arg, nil
}
