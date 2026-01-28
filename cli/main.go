package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nvlled/htmlformat"
)

func usage() {
	fmt.Printf("usage: %v [input-filename] [output-filename]\n", os.Args[0])
	fmt.Printf("  set - as the filename to read from stdin or write to stdout\n")
	os.Exit(0)
}

func main() {
	for _, arg := range os.Args {
		switch arg {
		case "-h", "--help":
			usage()
		}
	}

	var err error
	var input io.Reader
	var output io.Writer
	_ = input
	_ = output
	if len(os.Args) > 1 && os.Args[1] != "-" {
		inputFilename := os.Args[1]
		fmt.Printf("reading from file %s, ", os.Args[1])
		input, err = os.Open(inputFilename)
		if err != nil {
			fmt.Printf("failed to open file %s: %v\n", inputFilename, err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Print("reading from stdin, ")
		input = os.Stdin
	}

	if len(os.Args) > 2 && os.Args[2] != "-" {
		outputFilename := os.Args[2]
		fmt.Printf("writing to file %s\n", outputFilename)

		if os.Args[1] == outputFilename {
			fmt.Printf("input and output file must not be the same file\n")
			os.Exit(1)
		}

		output, err = os.Create(outputFilename)
		if err != nil {
			fmt.Printf("failed to open file %s: %v\n", outputFilename, err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("writing to stdout")
		output = os.Stdout
	}

	if len(os.Args) > 3 {
		fmt.Println("too many arguments, not up for debate")
		os.Exit(1)
	}

	htmlformat.Pipe(input, output)
}
