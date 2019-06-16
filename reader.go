package bug

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"unicode/utf8"
)

func init() {
	image.RegisterFormat("bug", string(rune(brailleCharOffset)), Decode, DecodeConfig)
	image.RegisterFormat("bug", string(rune(0x283f)), Decode, DecodeConfig)
}

// Decode creates a new BUG image from the given stream.
func Decode(r io.Reader) (image.Image, error) {
	return NewDecoder(r).Decode()
}

// Decoder handles the BUG decoding.
type Decoder struct {
	r io.Reader

	Threshold
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, Threshold: DefaultThreshold}
}

func (d *Decoder) WithThreshold(t Threshold) *Decoder {
	d.Threshold = t
	return d
}

func (d *Decoder) Decode() (image.Image, error) {
	// Consume the stream.
	buf, err := ioutil.ReadAll(d.r)
	if err != nil {
		return nil, err
	}

	// Split the lines.
	rows := bytes.Split(bytes.TrimSpace(buf), []byte{'\n'})
	if len(rows) == 0 {
		return nil, errors.New("empty BUG image")
	}

	// Count the width/height.
	width := utf8.RuneCount(rows[0])
	height := len(rows)

	// Create the new BUG image object.
	img := NewGray(image.Rectangle{
		Max: image.Point{
			X: width * 2,  // 2 cols per cell.
			Y: height * 4, // 4 rows per cell.
		},
	})
	img.Threshold = d.Threshold

	if len(rows) > len(img.content) {
		return nil, errors.New("invalid BUG image: corrupted column")
	}

	// Row by row.
	for row, line := range rows {
		cells := bytes.Runes(line)
		if len(cells) > len(img.content[row]) {
			return nil, errors.New("invalid BUG image: corrupted row")
		}
		// For each cell.
		for col, cell := range cells {
			// Remove the braillCharOffset to get the actual value.
			cellVal := uint8(cell - brailleCharOffset)

			// Store the value in the image object.
			img.content[row][col] = cellVal

			// Update the "real" image for each of the 8 pixel in the cell.
			x, y := col*2, row*4 // Pixel origin of the cell.
			for i := 0; i < 2; i++ {
				for j := 0; j < 4; j++ {
					// If set, use Opaque.
					if cellVal&unicodeOffset(x+i, y+j) == cellVal {
						img.Gray.Set(x+i, y+j, color.Opaque)
					} else {
						// Otherwise, transparent.
						img.Gray.Set(x+i, y+j, color.Transparent)
					}
				}
			}
		}
	}

	return img, nil
}

// DecodeConfig complies with image.RegisterFormat but is not used.
func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, errors.New("decodeConfig for BUG requires a full read, use image.Decode instead")
}
