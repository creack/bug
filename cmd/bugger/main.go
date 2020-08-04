package main

import (
	"flag"
	"image"
	"io"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/creack/bug"
)

// initFlags parses the cli input flags and validates them.
func initFlags() (verbose bool, inputPath, outputPath string) {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.StringVar(&inputPath, "in", "", "Path to the input image. Supports jpg/png.")
	flag.StringVar(&outputPath, "out", "", "Target BUG file path. If missing, prints to stdout.")

	flag.Parse()

	if inputPath == "" {
		log.Printf("Missing -in.")
		flag.Usage()
		os.Exit(1)
	}

	return verbose, inputPath, outputPath
}

func main() {
	// Init the flags.
	verbose, inputPath, outputPath := initFlags()

	// Load the input image.
	in, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Error opening the input file %q: %s.", inputPath, err)
	}
	// Decode it in memory.
	imgIn, format, err := image.Decode(in)
	if err != nil {
		log.Fatalf("Error decoding image file contents: %s.", err)
	}
	if verbose {
		log.Printf("Successfully decoded %q as %s.", inputPath, format)
	}

	// Convert the image
	imgOut := bug.Convert(imgIn, bug.DefaultThreshold.Inverse())

	// Create the target file if needed.
	var out io.WriteCloser
	if outputPath != "" {
		out, err = os.Create(outputPath)
		if err != nil {
			log.Fatalf("Error creating the output file %q: %s.", outputPath, err)
		}
		defer func() { _ = out.Close() }() // Best effort.
	} else {
		out = os.Stdout
	}

	if err := bug.Encode(out, imgOut); err != nil {
		log.Fatalf("Error encoding the result BUG image to the output file %q: %s.", outputPath, err)
	}
}
