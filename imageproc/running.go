package imageproc

import (
	"caddae/callout"
	"caddae/drawing"
	"caddae/types"
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// running handles proocessing of the running asbuilt image
func (ip *ImageProc) running() error {
	il := ip.log.With().Str("func", "running").Logger()
	il.Debug().Msg("starting running asbuilt process")

	var err error
	ip.ra.img, err = ip.OpenImage(ip.conf.Ra)
	if err != nil {
		il.Debug().Err(err).Msg("failed to open image")
		return errors.Wrapf(err, "ip.OpenImage(%s): failed to open image", ip.conf.Ra)
	}
	ip.ra.canvas.SetImage(ip.ra.img)

	ip.ra.cm, _, err = ip.preProcess(ip.ra.img, false)
	if err != nil {
		il.Debug().Err(err).Msg("failed to preprocess image")
		return errors.Wrapf(err, "failed to preprocess image")
	}
	il.Debug().Int("colorsFound", len(ip.ra.cm)).Send()

	il.Debug().Msg("Creating callout box")
	msg := fmt.Sprintf("Creating callout box and saving resulting image .. ")
	ip.UpdateUI(msg)

	prod := ip.CreateProdUnits()
	c := callout.New(prod, ip.ra.img)
	c.CreateCallout()
	c.AddCallout(ip.ra.img)

	msg = fmt.Sprintf("Drawing blue lines on running asbuilt ..")
	ip.UpdateUI(msg)

	//il.Debug().Interface("approxChanges", ip.ra.approxChanges)
	ip.ra.img = ip.ra.canvas.DrawLines(ip.ra.approxChanges, ip.RedlineEdge(), ip.RunningEdge())

	il.Debug().Msg("Saving updated running asbuilt")
	msg = "Saving updated running asbuilt file .."
	ip.UpdateUI(msg)

	f := ip.RunningFilePath()
	if err := ip.SaveRunning(f, "png"); err != nil {
		msg = fmt.Sprintf("ip.SaveUpdatedRunning(%s, %s): error saving updated redline file - %v", f, "png", err)
		ip.UpdateUI(msg)
	}

	msg = fmt.Sprintf("Running successfully saved as %s!\nEnd of application process. :)", f)
	ip.UpdateUI(msg)
	return nil
}

// Running returns the new running asbuilt image
func (ip *ImageProc) Running() image.Image {
	return ip.ra.img
}

// CreateProdUnits creates & returns the production for the running asbuilt
func (ip *ImageProc) CreateProdUnits() *types.Production {
	var p types.Production
	p.Date = ip.conf.Wpd

	var units []types.Unit
	if ip.conf.Strand != 0.0 {
		qty := fmt.Sprintf("%.2f", ip.conf.Strand)
		str := strings.Split(qty, ".")
		dec := str[1]
		if dec == "00" {
			qty = str[0]
		}
		name := "C300-01"
		txt := name + " = " + qty + "'"
		unit := types.Unit{
			Name:  name,
			Qty:   qty,
			Text:  txt,
			Color: drawing.Blue,
		}
		units = append(units, unit)
	}

	if ip.conf.Cable != 0.0 {
		qty := fmt.Sprintf("%.2f", ip.conf.Cable)
		str := strings.Split(qty, ".")
		dec := str[1]
		if dec == "00" {
			qty = str[0]
		}
		name := "C300-02"
		txt := name + " = " + qty + "'"
		unit := types.Unit{
			Name:  name,
			Qty:   qty,
			Text:  txt,
			Color: drawing.Blue,
		}
		units = append(units, unit)

		name2 := "C400"
		txt2 := name2 + " = " + qty + "'"
		unit2 := types.Unit{
			Name:  name2,
			Qty:   qty,
			Text:  txt2,
			Color: drawing.Blue,
		}
		units = append(units, unit2)
	}

	if ip.conf.Overlash != 0.0 {
		qty := fmt.Sprintf("%.2f", ip.conf.Overlash)
		str := strings.Split(qty, ".")
		dec := str[1]
		if dec == "00" {
			qty = str[0]
		}
		name := "C300-03"
		txt := name + " = " + qty + "'"
		unit := types.Unit{
			Name:  name,
			Qty:   qty,
			Text:  txt,
			Color: drawing.Blue,
		}
		units = append(units, unit)

		name2 := "C400"
		txt2 := name2 + " = " + qty + "'"
		unit2 := types.Unit{
			Name:  name2,
			Qty:   qty,
			Text:  txt2,
			Color: drawing.Blue,
		}
		units = append(units, unit2)
	}

	if ip.conf.Anchors != 0.0 {
		qty := fmt.Sprintf("%d", int(ip.conf.Anchors))
		name := "C300-04"
		txt := name + " = " + qty
		unit := types.Unit{
			Name:  name,
			Qty:   qty,
			Text:  txt,
			Color: drawing.Coral,
		}
		units = append(units, unit)
	}

	p.Units = units
	return &p
}

// RunningEdge returns the first black edge pixel found in the running asbuilt
//
// The reason we do this is because the scanned images are slightly off ceneterd
// from the original asbuilt image, so we'll use both images first black pixel
// on the right edge of the page, since we know all maps should have this.
//
// We'll use the difference between the two images x and y positions to
// determine how much we need to shift our approximate change pixels
func (ip *ImageProc) RunningEdge() *drawing.Pixel {
	var col color.RGBA
	var r, g, b, a uint32

	// Get the X, Y bounds of the image
	bnds := ip.ra.img.Bounds()
	var pixel drawing.Pixel

	// Now lets loop through it.
	for y := bnds.Max.Y; y > bnds.Min.Y; y-- {
		for x := bnds.Max.X; x > bnds.Min.X; x-- {
			r, g, b, a = ip.ra.img.At(x, y).RGBA()
			col.R, col.G, col.B, col.A = uint8(r), uint8(g), uint8(b), uint8(a)
			if v, ok := ip.ra.cm[col]; ok {
				if v.Range != nil && v.Range.Name == drawing.BLACKISH {
					return &drawing.Pixel{X: x, Y: y}
				}
			}
		}
	}
	return &pixel
}

// RunningFilePath gets the new file path for the updated running image
func (ip *ImageProc) RunningFilePath() string {
	// Sample now value: 2020-04-09 11:24:14.785868848 +0000 UTC m=+0.000187421
	now := time.Now()
	ts := now.Format(time.RFC3339)

	file := ip.conf.Ra
	filepath := strings.Split(file, "/")
	pathname := filepath[len(filepath)-1]
	path := file[:len(pathname)]
	f := path + "/caddae/edits/testfiles/" + ip.conf.Jn + "_" + ts + ".png"
	ip.ra.newFile = f
	return f
}
