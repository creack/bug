package main

import (
	"bufio"
	"bytes"
	"flag"
	"image"
	"io/ioutil"
	"log"
	"os"

	_ "image/png"

	"github.com/creack/bug"
)

func main() {
	var (
		verbose    = flag.Bool("v", false, "verbose mode")
		imageFile  = flag.String("img", "image.png", "path to image file")
		outputFile = flag.String("out", "", "path to output file (default stdout)")
	)
	flag.Parse()

	buf, err := ioutil.ReadFile(*imageFile)
	if err != nil {
		log.Fatalf("Error reading image file %q: %s", *imageFile, err)
	}
	img, format, err := image.Decode(bytes.NewBuffer(buf))
	if err != nil {
		log.Fatalf("Error decoding file contents: %s", err)
	}
	if *verbose {
		log.Printf("Successfully decoded %q as %s", *imageFile, format)
	}

	var out bytes.Buffer
	if err := bug.Encode(&out, img); err != nil {
		log.Fatalf("Error encoding image as bug: %s", err)
	}
	if *outputFile == "" {
		f := bufio.NewWriter(os.Stdout)
		f.Write(out.Bytes())
		f.Flush()
		return
	}
	if err := ioutil.WriteFile(*outputFile, out.Bytes(), 0644); err != nil {
		log.Fatalf("Error writing output to %q", *outputFile)
	}
}
