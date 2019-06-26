package main

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"

	"github.com/creack/bug"
	// Import the common image formats to make sure we don't
	// break anything there.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Make sure *Gray implements the image.Image interface.
var _ image.Image = (*bug.Gray)(nil)

// Load a bug file, decode, and then encode it.
func decodeEncodeBug(name string) {
	// Load the .bug file and decode it.
	img, _, err := image.Decode(mustGetFile("../testdata/" + name + ".bug"))
	if err != nil {
		log.Fatalf("Error decoding testdata image %q.", name)
	}
	// Re encode the .bug file.
	actual := bytes.NewBuffer(nil)
	if err := bug.Encode(actual, img); err != nil {
		log.Fatalf("Error encoding to bug: %s", err)
	}
}

// Convert PNG to bug.
func convertPNG(name string) {
	// Load png file and decode it.
	img, _, err := image.Decode(mustGetFile("../testdata/" + name + ".png"))
	if err != nil {
		log.Fatalf("Error decoding image: %s", err)
	}
	// Convert to .bug.
	converted := bug.Convert(img, bug.DefaultThreshold)
	// Encode the converted image in a buffer.
	actual := bytes.NewBuffer(nil)
	bug.Encode(actual, converted)
}

func main() {
	convertPNG("appenginegopher")
	decodeEncodeBug("appenginegopher")
}

func mustGetFile(path string) *bytes.Buffer {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %q: %s", path, err)
	}
	return bytes.NewBuffer(buf)
}
