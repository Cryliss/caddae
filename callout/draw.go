package callout

import (
	"caddae/drawing"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

// AddText adds text to the provided image.
func (c *Callout) AddText(img image.Image, x, y int, text string, col color.RGBA) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	d := &font.Drawer{
		Dst:  img.(*image.RGBA),
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

// HorizontalLine draws a horizontal line on the image.
func (c *Callout) HorizontalLine(x1, y, x2 int, img image.Image, col color.RGBA) {
	for ; x1 <= x2; x1++ {
		if canSet, ok := img.(drawing.CanSet); ok {
			canSet.Set(x1, y, col)
		}
	}
}

// VerticalLine draws a veritical line on the image.
func (c *Callout) VerticalLine(x, y1, y2 int, img image.Image, col color.RGBA) {
	for ; y1 <= y2; y1++ {
		if canSet, ok := img.(drawing.CanSet); ok {
			canSet.Set(x, y1, col)
		}
	}
}

// Rectangle draws a rectangle around the image.
func (c *Callout) Rectangle(x1, y1, x2, y2 int, img image.Image, border color.RGBA) {
	c.HorizontalLine(x1, y1, x2, img, border)
	c.HorizontalLine(x1, y2, x2, img, border)
	c.VerticalLine(x1, y1, y2, img, border)
	c.VerticalLine(x2, y1, y2, img, border)
}
