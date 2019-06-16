package bug

import (
	"image"
	"image/color"
)

// Using one braille unicode rune, we fit 2 cols and 4 rows, i.e.
// We fit 2 pixels by row and 4 pixels by columns.
// Neatly arranged so we can use binary operator to "merge" pixels in a cell.
var offsetMap = [4][2]uint8{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

// unicodeOffset returns the unicode offset for the given
// "real" pixel x, y.
func unicodeOffset(x, y int) uint8 {
	absX, absY := x%2, y%4
	if absX < 0 {
		absX += 2
	}
	if absY < 0 {
		absY += 4
	}
	return offsetMap[absY][absX]
}

// Braille chars start at 0x2800 (empty cell).
const brailleCharOffset rune = 0x2800

// Threshold to use in order to set a braille point in a cell.
// Negative value indivate inverse.
type Threshold int

const DefaultThreshold Threshold = 100

// Convert the given color to Opaque/Transparent.
func (cm Threshold) Convert(c color.Color) color.Color {
	switch c {
	case color.Opaque, color.White:
		c = color.Opaque
	case color.Transparent, color.Black:
		c = color.Transparent
	default:
		if color.GrayModel.Convert(c).(color.Gray).Y < uint8(cm) {
			c = color.Opaque
		} else {
			c = color.Transparent
		}
	}
	if cm >= 0 {
		return c
	}
	return color.Alpha16{A: ^c.(color.Alpha16).A}
}

// Inverse toggles the inverse mode for te threshold.
func (cm Threshold) Inverse() Threshold {
	return -cm
}

// Gray wraps a gray scale image with braille characters.
// Each braille character represents 2x4 actual pixels.
type Gray struct {
	// real holds the real pixel version of the image.
	*image.Gray

	// content holds the braille representation of the image.
	// Not using stdlib's single dim slice as benchmark shows
	// it is faster with 2 dim (i.e. without the mmath to map 1d to 2d).
	content [][]uint8

	// Rect is the image's bounds, in cells.
	Rect image.Rectangle

	// Threshold to toogle braille point based on gray scale.
	Threshold Threshold
}

// NewGray creates a new Black and White Braille Unicode Graphic (BUG) image.
// The rectangle is expected to be in "real" pixels.
func NewGray(r image.Rectangle) *Gray {
	width, height := r.Dx()/2, r.Dy()/4
	if r.Dx()%2 != 0 {
		width++
	}
	if r.Dy()%4 != 0 {
		height++
	}
	img := &Gray{
		Gray: image.NewGray(r),
		Rect: image.Rectangle{
			Min: r.Min,
			Max: image.Point{
				X: width,  // 2 cols per cell.
				Y: height, // 4 rows per cell.
			},
		},
		Threshold: DefaultThreshold,
	}
	img.content = make([][]uint8, img.Rect.Dy())
	for i := range img.content {
		img.content[i] = make([]uint8, img.Rect.Dx())
	}
	return img
}

// Clear all pixels.
func (p *Gray) Clear() {
	// The for loop is more efficient that a copy of empty element.
	// https://github.com/golang/go/issues/5373
	// See bench_content_test.go
	for i, l := range p.content {
		for j := range l {
			p.content[i][j] = 0
		}
	}
	for i := range p.Gray.Pix {
		p.Gray.Pix[i] = 0
	}
}

// BrailleAt returns the Braille Unicode for the given cell.
func (p *Gray) BrailleAt(col, row int) rune {
	// Discard pixels outside the image.
	if !(image.Point{col, row}.In(p.Rect)) {
		return brailleCharOffset
	}

	return rune(p.content[row][col]) + brailleCharOffset
}

// Set implements the image.Image interface.
// Update both the "real" version of the image, and the braille mapping.
func (p *Gray) Set(x, y int, c color.Color) {
	// Discard pixels outside the image.
	if !(image.Point{x, y}.In(p.Gray.Rect)) {
		return
	}
	p.Gray.Set(x, y, c)
	p.SetBraille(x, y, p.ColorModel().Convert(c))
}

// ColorModel implements the image.Image interface. It defines the
// grayscale threshold when to set the braille character point.
func (p *Gray) ColorModel() color.Model {
	return p.Threshold
}

// SetBraille updates the cell with the given "real" pixel x,y.
func (p *Gray) SetBraille(x, y int, c color.Color) {
	col, row := x/2, y/4

	// If opaque, set the point, otherwise, remove it.
	if c == color.Opaque {
		p.content[row][col] |= unicodeOffset(x, y)
	} else {
		p.content[row][col] &^= unicodeOffset(x, y)
	}
}
