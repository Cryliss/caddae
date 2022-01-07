package imageproc

import (
	"caddae/drawing"
	"caddae/types"
	"image"

	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog"
)

// Config for the running asbuilt
type Config struct {
	Rl       string
	Ra       string
	Jn       string
	Wpd      string
	Strand   float64
	Cable    float64
	Overlash float64
	Anchors  float64
}

// ImageProc data type for image processing
type ImageProc struct {
	log  zerolog.Logger
	conf Config
	rl   *Redline
	ra   *Running
	ui   bool
	u    types.UI
	g    *gocui.Gui
}

// Redline data type for the redline image
type Redline struct {
	newFile string
	img     image.Image
	cm      drawing.ColorMap
	bChange []*drawing.Pixel
	yChange []*drawing.Pixel
	wChange []*drawing.Pixel
}

// Running data type for the running asbuilt image
type Running struct {
	canvas        *drawing.Canvas
	newFile       string
	img           image.Image
	cm            drawing.ColorMap
	approxChanges []*drawing.Pixel
	bChange       []*drawing.Pixel
	yChange       []*drawing.Pixel
	wChange       []*drawing.Pixel
}
