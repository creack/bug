package bug

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"testing"

	// Import the common image formats to make sure we don't
	// break anything there.
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

// Make sure *Gray implements the image.Image interface.
var _ image.Image = (*Gray)(nil)

// Test decoding .bug file and re-encoding it.
func TestDecodeEncode(t *testing.T) {
	decodeEncode := func(t *testing.T, name string) {
		// Final expectation.
		expect := mustGetFile(t, "testdata/"+name+".bug")
		// Load the .bug file and decode it.
		img, _, err := image.Decode(mustGetFile(t, "testdata/"+name+".bug"))
		requireNoError(t, err, "Decode testdata image %q.", name)

		// Re encode the .bug file and assert it with original.
		actual := bytes.NewBuffer(nil)
		requireNoError(t, Encode(actual, img), "Encode decoded image %q.", name)
		assertEqual(t, expect, actual, "Unexpected re-encoded image.")
	}
	t.Run("appenginegopher", func(t *testing.T) { decodeEncode(t, "appenginegopher") })
	t.Run("biplane", func(t *testing.T) { decodeEncode(t, "biplane") })
}

// Test decoding .bug file and re-encoding it.
func TestDecodeEncodePNGDecodeEncode(t *testing.T) {
	decodeEncode := func(t *testing.T, name string) {
		// Final expectation.
		expect := mustGetFile(t, "testdata/"+name+".bug")
		// Load the .png file and decode it.
		img, _, err := image.Decode(mustGetFile(t, "testdata/"+name+".png"))
		requireNoError(t, err, "Decode testdata image %q.", name)

		// Encode the to .png.
		pngBuf := bytes.NewBuffer(nil)
		requireNoError(t, png.Encode(pngBuf, img), "Encode decoded image to png %q.", name)

		// Decode .png and convert to .bug.
		pngImg, _, err := image.Decode(pngBuf)
		requireNoError(t, err, "Decode testdata image %q.", name)

		// Convert to .bug.
		converted := Convert(pngImg, DefaultThreshold)
		// Encode the converted image in a buffer and assert.
		actual := bytes.NewBuffer(nil)
		requireNoError(t, Encode(actual, converted), "Encode converted image %q.", name)
		assertEqual(t, expect, actual, "Unexpected converted image.")

	}
	t.Run("appenginegopher", func(t *testing.T) { decodeEncode(t, "appenginegopher") })
	t.Run("biplane", func(t *testing.T) { decodeEncode(t, "biplane") })
}

// Test decoding .png file and converting it in inverse mode.
func TestDecodeEncodeInverse(t *testing.T) {
	decodeEncode := func(t *testing.T, name string) {
		// Final expectation.
		expect := mustGetFile(t, "testdata/"+name+".inverse.bug")
		// Load the .png file and decode it.
		img, _, err := image.Decode(mustGetFile(t, "testdata/"+name+".png"))
		requireNoError(t, err, "Decode testdata image %q.", name)

		// Conver to .bug in inverse mode.
		img = Convert(img, DefaultThreshold.Inverse())

		// Re encode the .bug file and assert it with original.
		actual := bytes.NewBuffer(nil)
		requireNoError(t, Encode(actual, img), "Encode converted image %q.", name)
		assertEqual(t, expect, actual, "Unexpected converted image.")
	}
	t.Run("appenginegopher", func(t *testing.T) { decodeEncode(t, "appenginegopher") })
	t.Run("biplane", func(t *testing.T) { decodeEncode(t, "biplane") })
}

// Test convert from .png to .bug.
func TestPNGConvert(t *testing.T) {
	convertImage := func(t *testing.T, name string) {
		// Final expectation.
		expect := mustGetFile(t, "testdata/"+name+".bug")
		// Load png file and decode it.
		img, _, err := image.Decode(mustGetFile(t, "testdata/"+name+".png"))
		requireNoError(t, err, "Decode testdata image %q.", name)
		// Convert to .bug.
		converted := Convert(img, DefaultThreshold)
		// Encode the converted image in a buffer and assert.
		actual := bytes.NewBuffer(nil)
		requireNoError(t, Encode(actual, converted), "Encode converted image %q.", name)
		assertEqual(t, expect, actual, "Unexpected converted image.")
	}
	t.Run("appenginegopher", func(t *testing.T) { convertImage(t, "appenginegopher") })
	t.Run("biplane", func(t *testing.T) { convertImage(t, "biplane") })
}

// Make sure the stdlib formats are still working.
func TestStdlibFormats(t *testing.T) {
	loadImage := func(t *testing.T, pth, typ string) {
		t.Helper()

		_, name, err := image.Decode(mustGetFile(t, pth))
		requireNoError(t, err, "Decode testdata file %q (%q).", pth, typ)
		assertEqual(t, typ, name, "Unexpected image type name for testdata file %q (%q).", pth, typ)
	}

	t.Run("png", func(t *testing.T) { loadImage(t, "testdata/video-001.png", "png") })
	t.Run("gif", func(t *testing.T) { loadImage(t, "testdata/video-001.gif", "gif") })
	t.Run("jpeg", func(t *testing.T) { loadImage(t, "testdata/video-001.jpeg", "jpeg") })
}

func assertEqual(tb testing.TB, expect, actual interface{}, msg string, args ...interface{}) bool {
	tb.Helper()

	if expect, actual := fmt.Sprint(expect), fmt.Sprint(actual); expect != actual {
		tb.Errorf("Assert failed.\n%s\nExpect:\t%s\nActual:\t%s", fmt.Sprintf(msg, args...), expect, actual)
		return false
	}
	return true
}

func requireNoError(tb testing.TB, err error, msg string, args ...interface{}) {
	tb.Helper()

	if err != nil {
		tb.Fatalf("Unexpected error: %s\n%s", err, fmt.Sprintf(msg, args...))
	}
}

func mustGetFile(tb testing.TB, pth string) *bytes.Buffer {
	tb.Helper()
	buf, err := ioutil.ReadFile(pth)
	requireNoError(tb, err, "Read file %q.", pth)
	return bytes.NewBuffer(buf)
}
