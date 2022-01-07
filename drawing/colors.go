package drawing

import (
	"image"
	"image/color"
)

// GetColors gets (and if change is set, changes) the colors in the provided image.
//
// Note that this is *not* efficent, as a better method would be to get the actual image.Image type and iterate
// over the pixels directly, saving a whole lot of conversion.
//
// Why? This was done quickly more as a proof of concept, so performance is not a consideration.
// Also this assumes the image provided is no higher then 24-bit color depth, and will happily fail if otherwise.
//
// This is simply because we "convert" the 32-bit RGB values to 8-bit values without any real checking or conversion to account for it.
// Probably should use RGBA64, which would handle 48-bit color depth, if you actually want to do something besides a proof of concept here.
//
// I believe GIF only supports a max of 24-bit color depth, where as PNG and JPEG XT supports 48-bit.
//
// For more details, see - https://en.wikipedia.org/wiki/Color_depth
//
// This should be considered a known bug.
//
// Again, more proof of concept then anything I'd actually recommend using in production.
func (c *Canvas) GetColors(in image.Image, change bool) (ColorMap, ChangeMap) {
	var col color.RGBA
	var r, g, b, a uint32
	cm := make(ColorMap, 1)
	chm := make(ChangeMap, 3)

	var bPixels, wPixels, yPixels []*Pixel

	// Get the X, Y bounds of the image
	bnds := in.Bounds()

	// Now lets loop through it.
	for x := bnds.Min.X; x < bnds.Max.X; x++ {
		for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
			// This? Yeah, this is a performance killer, all the conversion done in these functions.
			r, g, b, a = in.At(x, y).RGBA()

			// And here? This is a bug waiting to bite some unsuspecting image.
			//
			// We just chop off 3 bytes from the colors and hope they were really 8-bit values to begin with.
			col.R, col.G, col.B, col.A = uint8(r), uint8(g), uint8(b), uint8(a)

			// fmt.Printf("At %d x %d, color %#v\n", x, y, col)

			if v, ok := cm[col]; ok {
				// This color has already been looked up.
				//
				// First, update the count.
				cm[col] = &ColorCount{v.Count + 1, v.Range}

				// Now we need to check if it has a range, and if that range wants us to change the color of the
				// current pixel or not?
				if change && v.Range != nil && v.Range.Replace {
					// Yep, change the pixel.
					//
					// Note that image.Image interface doesn't have a Set() function, but the types that typically make it up do.
					//
					// So we see if it supports the Set() function via our canSetImage interface, which lets us know and provides us
					// the interface to call Set() with.
					//
					// Note this will fail for JPEGs, as those are typically image.YCbCr, which does not have a Set() function.
					if canSet, ok := in.(CanSet); ok {
						canSet.Set(x, y, v.Range.Make)

						/*fmt.Printf("At %d x %d, changed pixel from #%02x%02x%02x to #%02x%02x%02x\n", x, y,
						col.R, col.G, col.B, v.Range.Make.R, v.Range.Make.G, v.Range.Make.B )*/
					}
					pc := Pixel{
						X: x,
						Y: y,
					}
					if v.Range.Name == YELLOWISH {
						yPixels = append(yPixels, &pc)
					} else if v.Range.Name == BLACKISH {
						bPixels = append(bPixels, &pc)
					} else {
						wPixels = append(wPixels, &pc)
					}
				}
			} else {
				//fmt.Printf("First seen color: #%02x%02x%02x, %x\n", col.R, col.G, col.B, col.A)
				cm[col] = &ColorCount{1, c.GetRange(col)}
			}
		}
	}
	chm[BLACKISH] = bPixels
	chm[WHITEISH] = wPixels
	chm[YELLOWISH] = yPixels
	return cm, chm
}

// ChangeColors changes the yellow pixels in the image to blue.
func (c *Canvas) ChangeColors(cm ColorMap, in image.Image, rName string) (ColorMap, ChangeMap) {
	var col color.RGBA
	var r, g, b, a uint32
	chm := make(ChangeMap, 1)

	var bPixels, wPixels, yPixels []*Pixel

	// Get the X, Y bounds of the image
	bnds := in.Bounds()

	// Now lets loop through it.
	for x := bnds.Min.X; x < bnds.Max.X; x++ {
		for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
			r, g, b, a = in.At(x, y).RGBA()
			col.R, col.G, col.B, col.A = uint8(r), uint8(g), uint8(b), uint8(a)

			// fmt.Printf("At %d x %d, color %#v\n", x, y, col)

			if v, ok := cm[col]; ok {
				if v.Range != nil && v.Range.Name == rName {
					// Yep, change the pixel.
					if canSet, ok := in.(CanSet); ok {
						canSet.Set(x, y, v.Range.Make)
					}
					pc := Pixel{
						X: x,
						Y: y,
					}

					if v.Range.Name == YELLOWISH {
						yPixels = append(yPixels, &pc)
					} else if v.Range.Name == BLACKISH {
						bPixels = append(bPixels, &pc)
					} else {
						wPixels = append(wPixels, &pc)
					}
				}
			} else {
				//fmt.Printf("First seen color: #%02x%02x%02x, %x\n", col.R, col.G, col.B, col.A)
				cm[col] = &ColorCount{1, c.GetRange(col)}
			}
		}
	}
	chm[BLACKISH] = bPixels
	chm[WHITEISH] = wPixels
	chm[YELLOWISH] = yPixels
	return cm, chm
}

// GetRange checks if a range is found for this specific color, return it.
func (c *Canvas) GetRange(in color.RGBA) *ColorRange {
	// Lets see if the provided color has a range or not.
	for _, cr := range Ranges {
		// red match
		if in.R > cr.RMax || in.R <= cr.RMin {
			continue
		}

		if in.G > cr.GMax || in.G <= cr.GMin {
			continue
		}

		if in.B > cr.BMax || in.B <= cr.BMin {
			continue
		}

		// If we are here, that means the range matches, so return it.
		return &cr
	}

	return nil
}
