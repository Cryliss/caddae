// Package callout creates production callouts for the running asbuilt.
package callout

import (
	"caddae/drawing"
	"caddae/types"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

// New creates and returns a new callout
func New(prod *types.Production, img image.Image) *Callout {
	// Initialize the callout.
	var c Callout
	c.prod = prod

	// Resize the callout based on image dimensions
	bnds := img.Bounds()
	c.Resize(bnds.Max.X, bnds.Max.Y)

	// Determine how many production units we have to see if we need to resize.
	numUnits := len(prod.Units)
	if numUnits > 1 {
		numUnits--
		addtnlHeight := numUnits * addtlnProdHeight
		c.dim.y2 = c.dim.y2 + addtnlHeight
		canvasHeight = canvasHeight + addtnlHeight
	}

	c.canvas = image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	return &c
}

// Resize resizes the callout and production canvas dimensions and text dimensions
// based on the provides image bounds
func (c *Callout) Resize(xMax, yMax int) {
	cw := math.Round(float64(xMax) * 0.05)
	ch := math.Round(float64(yMax) * 0.05)
	canvasWidth = int(cw)
	canvasHeight = int(ch)

	// Create a Dimension for the dimensions of the callout.
	var dim Dimensions
	dim.x1 = 0
	dim.y1 = 0
	dim.x2 = canvasWidth - 1
	dim.y2 = canvasHeight - 1

	// Create a Text for the location of the date.
	var txt Text
	x := math.Round(float64(canvasWidth) * 0.3)
	y := math.Round(float64(canvasHeight) * 0.28125)

	txt.x = int(x)
	txt.y = int(y)

	c.dim = dim
	c.date = txt

	pw := math.Round(float64(canvasWidth) * 0.875)
	ph := math.Round(float64(canvasHeight) * 0.333)

	prodWidth = int(pw)
	prodHeight = int(ph)

	ap := math.Round(float64(prodHeight) * 1.15)
	ag := math.Round(float64(prodHeight) * 1.25)
	addtlnProdHeight = int(ap)
	prodGap = int(ag)

	px := math.Round(float64(canvasWidth) * 0.0625)
	py := math.Round(float64(canvasHeight) * 0.4167)

	prodDims.x1 = int(px)
	prodDims.y1 = int(py)
	prodDims.x2 = prodDims.x1 + prodWidth
	prodDims.y2 = prodDims.y1 + prodHeight

	ptx := math.Round(float64(prodDims.x1) * 1.1429)
	pty := math.Round(float64(prodDims.y1) * 1.525)
	prodText.x = int(ptx)
	prodText.y = int(pty)
}

// CreateCallout creates the callout canvas, adds the date and production boxes to the canvas.
func (c *Callout) CreateCallout() error {
	// Make a white mask over the callout image.
	draw.DrawMask(c.canvas.(draw.Image), image.Rect(c.dim.x1, c.dim.y1, c.dim.x2, c.dim.y2), &image.Uniform{drawing.White}, image.ZP, nil, image.ZP, draw.Src)

	// Add a rectange outlining the canvas.
	c.Rectangle(c.dim.x1, c.dim.y1, c.dim.x2, c.dim.y2, c.canvas, drawing.Black)
	c.Rectangle(c.dim.x1+1, c.dim.y1+1, c.dim.x2-1, c.dim.y2-1, c.canvas, drawing.Black)

	// Add the date to the canvas.
	c.AddText(c.canvas, c.date.x, c.date.y, c.prod.Date, drawing.Black)

	// For each production, create a prodbox.
	for _, u := range c.prod.Units {
		c.AddProdBox(u.Text, u.Color)
	}

	// Save our drawing of the canvas for review.
	c.SaveDrawing("/Users/sabra/go/src/caddae/edits/draw.png", c.canvas)
	return nil
}

// AddCallout adds the callout to the running asbuilt image.
func (c *Callout) AddCallout(img image.Image) {
	// Get the bounds of the image
	bnds := img.Bounds()

	// Get the max x & y values
	xMax := bnds.Max.X
	yMax := bnds.Max.Y

	// Get the middle x & y values of the image
	// (We're just gonna place it smack dab in the middle of the image for
	// the time being lol.)
	x := int(xMax/2) + int(xMax/4)
	y := int(yMax / 2)

	// Draw the callout onto the image
	draw.DrawMask(img.(draw.Image), image.Rect(x, y, xMax, yMax), c.canvas.(draw.Image), image.ZP, nil, image.ZP, draw.Src)
}

// AddProdBox adds a new production box to the callout
//
// For each new prod box, shift y2 down by 25.
// If more than 6 prod boxes, need to resize canvas
// Each prod box needs it's texts y value shifted down 30
// Each prod box need it's y1 & y2 values shifted down 30
func (c *Callout) AddProdBox(text string, col color.Color) {
	// Create a new canvas image for the production box
	canvas := image.NewRGBA(image.Rect(0, 0, prodWidth, prodHeight))

	// Do we already have a prod box?
	if numProd > 0 {
		// Yes, so let's shift this one down
		prodDims.y1 = prodDims.y1 + prodGap
		prodDims.y2 = prodDims.y2 + prodGap
		prodText.y = prodText.y + prodGap
	}

	// Now let's go ahead and position the text inside the box
	if len(text) <= 11 {
		addWidth := 10 - len(text)
		prodText.x = prodText.x + addWidth
	} else if len(text) < 15 {
		addWidth := 25 - len(text)
		prodText.x = prodText.x + addWidth
	} else if len(text) < 17 {
		addWidth := 35 - len(text)
		prodText.x = prodText.x + addWidth
	} else if len(text) < 25 {
		addWidth := 30 - len(text)
		prodText.x = prodText.x + addWidth
	} else {
		math.Round(float64(prodDims.x1) * 1.1429)
	}

	// Create a blue mask over the canvas
	draw.DrawMask(canvas, canvas.Bounds(), &image.Uniform{col}, image.ZP, nil, image.ZP, draw.Src)

	// Now, add our blue box to the callout canvas
	draw.DrawMask(c.canvas.(draw.Image), image.Rect(prodDims.x1, prodDims.y1, prodDims.x2, prodDims.y2), canvas, image.ZP, nil, image.ZP, draw.Src)

	// Draw a rectangle around the blue canvas
	c.Rectangle(prodDims.x1, prodDims.y1, prodDims.x2, prodDims.y2, c.canvas, drawing.Black)
	c.Rectangle(prodDims.x1+1, prodDims.y1+1, prodDims.x2-1, prodDims.y2-1, c.canvas, drawing.Black)

	// Add the production boxes text to the canvas
	c.AddText(c.canvas, prodText.x, prodText.y, text, drawing.Black)
	c.AddText(c.canvas, prodText.x, prodText.y, text, drawing.Black)
	// Increment the amount of production boxes we have
	numProd++
}

// SaveDrawing saves our callout drawing as a png.
func (c *Callout) SaveDrawing(out string, img image.Image) error {
	// We first create a temporary file, then if everything is OK we rename it.
	// This ensures we don't replace the output with any half-written files that
	// could break anything further down the line trying to read our output.
	newFile := out + ".tmp"

	f, err := os.OpenFile(newFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		e := fmt.Sprintf("Error writing file (%s): %s\n", newFile, err)
		return errors.New(e)
	}

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		e := fmt.Sprintf("Error encoding file (%s): %s\n", newFile, err)
		return errors.New(e)
	}

	// Ensure the contents are actually written to disk before we do the rename
	if err := f.Sync(); err != nil {
		e := fmt.Sprintf("sync(%s): %s", newFile, err)
		return errors.New(e)
	}

	f.Close()

	// Now rename the output.
	if err := os.Rename(newFile, out); err != nil {
		e := fmt.Sprintf("rename(%s, %s): %s", newFile, out, err)
		return errors.New(e)
	}

	return nil
}
