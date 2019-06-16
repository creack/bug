package bug

import (
	"image"
	"image/draw"
	"io"
	"unicode/utf8"
)

// Encode the BUG image.
func Encode(w io.Writer, img image.Image) error {
	bugImg := Convert(img)
	line := make([]byte, bugImg.Rect.Dx()*3+1) // 3 bytes per braille rune. + 1 for the newline.
	line[bugImg.Rect.Dx()*3] = '\n'
	for _, row := range bugImg.content {
		for i, cell := range row {
			utf8.EncodeRune(line[i*3:], rune(cell)+brailleCharOffset)
		}
		if _, err := w.Write(line); err != nil {
			return err
		}
	}
	return nil
}

// Convert the given image to a grayscale BUG one.
func Convert(img image.Image) *Gray {
	if g, ok := img.(*Gray); ok {
		return g
	}

	g := NewGray(img.Bounds())
	draw.Draw(g, g.Bounds(), img, img.Bounds().Min, draw.Over)
	return g
}
