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
func initFlags() (threshold int, inputPath, outputPath string) {
	flag.IntVar(&threshold, "t", 100, "Threshold for conversion. Set to negative for inverse output.")
	flag.StringVar(&inputPath, "in", "", "Path to the input image. Supports jpg/png.")
	flag.StringVar(&outputPath, "out", "", "Target BUG file path. If missing, prints to stdout.")

	flag.Parse()

	if inputPath == "" {
		log.Printf("Missing -in.")
		flag.Usage()
		os.Exit(1)
	}

	return threshold, inputPath, outputPath
}

func main() {
	// Init the flags.
	threshold, inputPath, outputPath := initFlags()

	// Load the input image.
	in, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Error opening the input file %q: %s.", inputPath, err)
	}
	// Decode it in memory.
	imgIn, _, err := image.Decode(in)
	if err != nil {
		log.Fatalf("Error decoding image file contents: %s.", err)
	}

	// Convert the image
	imgOut := bug.Convert(imgIn, bug.Threshold(threshold))

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
