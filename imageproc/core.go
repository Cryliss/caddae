package imageproc

import (
	"caddae/drawing"
	"caddae/types"
	"fmt"
	"image"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog"
)

// New initializes and returns a new imageproc
func New(conf Config, logger *zerolog.Logger) *ImageProc {
	var rl Redline
	var ra Running

	i := ImageProc{
		conf: conf,
		log:  logger.With().Str("module", "imageproc").Logger(),
		rl:   &rl,
		ra:   &ra,
	}
	il := i.log.With().Str("func", "New").Logger()
	il.Debug().Msg("Created")

	ra.canvas = drawing.New(&i.log)
	return &i
}

// ProcessImages starts the image processing for the redline and running asbuilt
func (ip *ImageProc) ProcessImages(u types.UI, g *gocui.Gui) error {
	il := ip.log.With().Str("func", "ProcessImages").Logger()
	il.Debug().Msg("Processing redline image")

	if u == nil && g == nil {
		il.Debug().Msg("UI logging is not set")
		ip.ui = false
	} else {
		il.Debug().Msg("UI logging is set")
		ip.ui = true
	}
	ip.u = u
	ip.g = g

	if err := ip.redline(); err != nil {
		return err
	}
	il.Debug().Msg("Processing running image")
	return ip.running()
}

// preProcess gets and changes the blackish and whiteish colors in the image
func (ip *ImageProc) preProcess(img image.Image, redline bool) (drawing.ColorMap, drawing.ChangeMap, error) {
	il := ip.log.With().Str("func", "preProcess").Logger()
	il.Debug().Bool("redline", redline).Send()

	// We went to get the colors in the image & change any that are "whiteish"
	// or blackish while we're at it, so let's just make sure those ranges
	// are to set to replace before we call it.
	for _, r := range drawing.Ranges {
		if r.Name == drawing.WHITEISH || r.Name == drawing.BLACKISH {
			r.Replace = true
		} else {
			r.Replace = false
		}
	}

	// Now lets get the colors in the image.
	cm, chm := ip.ra.canvas.GetColors(img, true)
	return cm, chm, nil
}

// preProcessColors gets and changes the yellowish colors in the image
func (ip *ImageProc) preProcessColors(img image.Image, redline bool) (drawing.ColorMap, error) {
	il := ip.log.With().Str("func", "preProcessColors").Logger()
	il.Debug().Bool("redline", redline).Msg("preProcessColors")

	// We went to get the colors in the image & change any that are "whiteish"
	// or blackish while we're at it, so let's just make sure those ranges
	// are to set to replace before we call it.
	for _, r := range drawing.Ranges {
		if r.Name == drawing.WHITEISH || r.Name == drawing.BLACKISH {
			r.Replace = false
		} else {
			r.Replace = true
		}
	}

	// Now lets get the colors in the image.
	cm, _ := ip.ra.canvas.GetColors(img, true)

	return cm, nil
}

// UpdateUI is a useful function to call anytime we want to give the user a
// messgae on the application log
func (ip *ImageProc) UpdateUI(msg string) {
	if !ip.ui {
		fmt.Fprintf(os.Stdout, msg)
		return
	}

	ip.g.Update(func(*gocui.Gui) error {
		if err := ip.u.Log(msg); err != nil {
			return err
		}
		return nil
	})
}
