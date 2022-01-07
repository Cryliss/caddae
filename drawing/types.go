package drawing

import (
	"github.com/rs/zerolog"
	"image"
	"image/color"
)

// Canvas to draw on
type Canvas struct {
	img image.Image
	log zerolog.Logger
}

// Pixel is an x,y point on the image
type Pixel struct {
	X, Y int
}

// Line is an array of pixels that we'll create from the approximate changes
type Line []*Pixel

// Lines is a nicer way of declaring an array of line objects
type Lines []Line

// CanSet is our quick interface that allows us to use the Set() function from the various types used to make an image.Image.
type CanSet interface {
	Set(x, y int, c color.Color)
}

// ColorCount is the counts of each range
type ColorCount struct {
	Count uint32
	Range *ColorRange
}

// ChangeMap is the map of changes we made by named range
type ChangeMap map[string][]*Pixel

// ColorMap is a map of colors to their color counts
type ColorMap map[color.RGBA]*ColorCount

// ColorItem is for counting each color item
type ColorItem struct {
	Count uint32
	Color color.RGBA
}

// ColorList is a nice data type for an array of color items
type ColorList []ColorItem

// Len returns the length of the ColorList
func (c ColorList) Len() int { return len(c) }

// Swap swaps two values in the color list array
func (c ColorList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less checks if the count of one color is less than another color in the list
func (c ColorList) Less(i, j int) bool { return c[i].Count > c[j].Count }

// RangeItem is an item in a range
type RangeItem struct{}

// RangeList is a nicer way of declaring an array of RangeItems
type RangeList []RangeItem

// ColorRange holds min/max values for colors
type ColorRange struct {
	// The R, G, B minimum and maximum range to
	// match this color.
	RMin uint8
	RMax uint8
	GMin uint8
	GMax uint8
	BMin uint8
	BMax uint8

	// The name given to this range
	Name string

	// If we replace this color with another
	Replace bool

	// The RGBA to replace the input color with
	Make color.RGBA
}

// ColorRanges is a nicer way of declaring an array of ColorRange
type ColorRanges []ColorRange
