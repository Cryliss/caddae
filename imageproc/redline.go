package imageproc

import (
	"caddae/drawing"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/pkg/errors"
)

// redline handles processing the redline image.
func (ip *ImageProc) redline() error {
	il := ip.log.With().Str("func", "redline").Logger()
	var err error

	// Open the image
	ip.rl.img, err = ip.OpenImage(ip.conf.Rl)
	if err != nil {
		return errors.Wrapf(err, "ip.OpenImage(%s): failed to open image", ip.conf.Rl)
	}

	var chm drawing.ChangeMap

	// Preprocess the image (change the "whiteish" colors to white, "blackish" colors to black)
	ip.rl.cm, _, err = ip.preProcess(ip.rl.img, true)
	if err != nil {
		return err
	}
	il.Debug().Int("colorsFound", len(ip.rl.cm)).Send()

	il.Debug().Msg("Changing yellow pixels to blue")
	_, err = ip.preProcessColors(ip.rl.img, true)
	if err != nil {
		return err
	}

	// Change the yellow pixels to blue
	_, chm = ip.ra.canvas.ChangeColors(ip.rl.cm, ip.rl.img, drawing.YELLOWISH)

	// Set the running approximate pixel changes equal to the redlines yellow changes
	ip.rl.yChange = chm[drawing.YELLOWISH]
	ip.ra.approxChanges = chm[drawing.YELLOWISH]
	//il.Debug().Interface("approxChanges", ip.ra.approxChanges).Send()

	msg := "Saving updated redline file .. \n"
	ip.UpdateUI(msg)

	f := ip.RedlineFilePath()
	il.Debug().Str("newRlFile", f)
	if err := ip.SaveRedline(f, "png"); err != nil {
		il.Debug().Err(err).Msg("failed to save updated redline file")

		msg = fmt.Sprintf("ip.SaveUpdatedRedline(%s, %s): error saving updated redline file - %v", f, "png", err)
		ip.UpdateUI(msg)
	}

	msg = fmt.Sprintf("Redline successfully saved as %s!\n", f)
	ip.UpdateUI(msg)
	return nil
}

// Redline returns the redline image
func (ip *ImageProc) Redline() image.Image {
	return ip.rl.img
}

// RedlineEdge returns the first black edge pixel found in the redline
// starting from the right edge of the image
func (ip *ImageProc) RedlineEdge() *drawing.Pixel {
	var col color.RGBA
	var r, g, b, a uint32

	// Get the X, Y bounds of the image
	bnds := ip.rl.img.Bounds()

	// Now lets loop through it.
	for y := bnds.Max.Y; y > bnds.Min.Y; y-- {
		for x := bnds.Max.X; x > bnds.Min.X; x-- {
			r, g, b, a = ip.rl.img.At(x, y).RGBA()
			col.R, col.G, col.B, col.A = uint8(r), uint8(g), uint8(b), uint8(a)
			if v, ok := ip.rl.cm[col]; ok {
				if v.Range != nil && v.Range.Name == drawing.BLACKISH {
					return &drawing.Pixel{X: x, Y: y}
				}
			}
		}
	}
	return &drawing.Pixel{}
}

// RedlineFilePath returns the new file path for the updated redline image
func (ip *ImageProc) RedlineFilePath() string {
	file := ip.conf.Rl

	filepath := strings.Split(file, "/")
	pathname := filepath[len(filepath)-1]

	path := file[:len(pathname)-1]
	fname := strings.Split(file, ".")
	name := fname[0]

	n := name[len(pathname)-1:]
	f := path + "edits/" + n + "_PREPROCESS.png"
	ip.rl.newFile = f
	return f
}
