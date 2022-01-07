package callout

import (
	"caddae/types"
	"image"
)

// Width and height of our callout canvas.
var canvasWidth int = 161
var canvasHeight int = 61

// Additional height we need on the callout canvas per prod box.
var addtlnProdHeight int = 20

// Gap between the starting y values of each prod box.
var prodGap int = 20

// Width and height of the production box.
var prodWidth int = 140
var prodHeight int = 20

// Starting dimensions of the production boxes.
var prodDims = Dimensions{
	x1: 10,
	y1: 25,
	x2: 150,
	y2: 45,
}

// Starting x & y values of the production text.
var prodText = Text{
	x: 20,
	y: 40,
}

// Number of production boxes currently made.
var numProd int = 0

// Callout that lists the work that was done that day, using production units
type Callout struct {
	date   Text
	prod   *types.Production
	dim    Dimensions
	canvas image.Image
	box    image.Image
	arrow  image.Image
}

// Dimenstions of each box
type Dimensions struct {
	x1, y1 int
	x2, y2 int
}

// Text is the location of text in the boxes.
type Text struct {
	x, y int
}
