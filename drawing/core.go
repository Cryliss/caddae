// Package drawing provides functionality for drawing on an image
package drawing

import (
	"image"
	"image/color"
	"math"

	"github.com/rs/zerolog"
)

// New returns a new Canvas
func New(l *zerolog.Logger) *Canvas {
	var c Canvas

	c.log = l.With().Str("module", "canvas").Logger()
	cl := c.log.With().Str("func", "New").Logger()
	cl.Debug().Msg("Created")

	return &c
}

// SetImage sets the image the canvas will be using
func (c *Canvas) SetImage(img image.Image) {
	c.img = img
}

// DrawLines takes the approximate changes retrieved from the redline and
// shifts them closer to the correct location, splits them into separate
// straight lines, and then draws an antialiaed line
func (c *Canvas) DrawLines(approxChanges []*Pixel, firstRlBlack, firstRaBlack *Pixel) image.Image {
	cl := c.log.With().Str("func", "DrawLines").Logger()
	cl.Debug().Interface("firstRlBlack", firstRlBlack).Msg("First black pixel in the redline")
	cl.Debug().Interface("firstRaBlack", firstRaBlack).Msg("First black pixel in the running")

	difX := firstRaBlack.X - firstRlBlack.X
	difY := firstRaBlack.Y - firstRlBlack.Y

	// Shift the pixels
	shifted := c.ShiftPixels(approxChanges, difX, difY)

	lines := c.ConvertLines(shifted)
	for i, line := range lines {
		cl.Debug().Interface("line", line).Msg("next line")
		if i > 5 {
			c.DrawAntialiased(*line[0], *line[len(line)-1], Blue)
		}
	}
	return c.img
}

// GetColor gets the color of the pixel at (x,y)
func (c *Canvas) GetColor(x, y int) color.RGBA {
	return c.img.At(x, y).(color.RGBA)
}

var avgX, avgY, stdY, stdX float64
var xPlus, xMinus, yPlus, yMinus float64

// ShiftPixels shifts the pixels in the approximate changes
func (c *Canvas) ShiftPixels(approxChanges []*Pixel, shiftX, shiftY int) []*Pixel {
	cl := c.log.With().Str("func", "ShiftPixels").Logger()
	var pixels []*Pixel

	bnds := c.img.Bounds()
	xMax := bnds.Max.X
	yMax := bnds.Max.Y
	shiftX = 60
	shiftY = 45

	var sumX, sumY int

	for _, pixel := range approxChanges {
		newX := math.Min(float64(pixel.X+shiftX), float64(xMax))
		newY := math.Min(float64(pixel.Y-shiftY), float64(yMax))
		x,y := c.GetNearestBlack(int(newX), int(newY))
		newP := &Pixel{X: x, Y: y}
		pixels = append(pixels, newP)

		sumX += pixel.X
		sumY += pixel.Y
	}

	//cl.Debug().Interface("newPixels", pixels).Msg("Pixels after shift")

	c.analysis(sumX, sumY, pixels)
	analysis := map[string]float64{
		"sumX":         float64(sumX),
		"sumY":         float64(sumY),
		"averageX":     avgX,
		"averageY":     avgY,
		"standardDevX": stdX,
		"standardDevY": stdY,
		"xPlus":        xPlus,
		"yPlus":        yPlus,
		"xMinus":       xMinus,
		"yMinus":       yMinus,
	}

	cl.Debug().Interface("analysis", analysis).Msg("Analysis on pixels")
	//ra.saveApproxChanges(pixels)
	return pixels
}

// analysis performs statistical analysis on the approxChange values
func (c *Canvas) analysis(sumX, sumY int, pixels []*Pixel) {
	n := len(pixels)
	if n == 0 {
		return
	}
	avgX, avgY = float64(sumX/n), float64(sumY/n)
	c.stdDev(pixels)

	xPlus, xMinus = avgX+stdX, avgX-stdX
	yPlus, yMinus = avgY+stdY, avgY-stdY
}

// stdDev calculates the standard deviation
func (c *Canvas) stdDev(pixels []*Pixel) {
	var sumX, sumY float64
	for _, p := range pixels {
		x, y := float64(p.X), float64(p.Y)
		sumX += math.Abs(x-avgX) * math.Abs(x-avgX)
		sumY += math.Abs(y-avgY) * math.Abs(y-avgY)
	}
	n := float64(len(pixels))
	stdX, stdY = math.Sqrt(sumX/n), math.Sqrt(sumY/n)
}

// ConvertLines converts an array of pixels to individual lines
func (c *Canvas) ConvertLines(pixels []*Pixel) Lines {
	cl := c.log.With().Str("func", "ConvertLines").Logger()
	cl.Debug().Msg("Started")

	var lines Lines
	lines = make(Lines, 1)

	var line Line
	var weight float64 = 0.5

	for _, pixel := range pixels {
		x, y := c.GetNearestBlack(pixel.X, pixel.Y)
		if x == -1 {
			continue
		}

		// Create a new pixel at the location on the canvas
		cPixel := &Pixel{X: x, Y: y}

		// If we have no pixels in the line, let's just add it to the line
		if len(line) == 0 || len(line) == 1 {
			line = append(line, cPixel)
			continue
		}

		// We have a pixel in the line, but do we have any other lines?
		if len(lines) == 0 {
			// No, we don't, so let's check if this pixel is on the current line
			if c.OnLine(line, cPixel, &weight) {
				cl.Debug().Interface("pixel", cPixel).Msg("on line lines == 0")
				// It is, so let's add it to the line and continue
				line = append(line, cPixel)
				continue
			}
			cl.Debug().Interface("pixel", cPixel).Msg("not on line lines == 0")
			// It wasn't on the line, so let's add the line to lines and create
			// a new line to append the pixel to
			lines = append(lines, line)

			var newline Line
			line = newline
			line = append(line, cPixel)
			continue
		}

		if c.OnLine(line, cPixel, &weight) {
			cl.Debug().Interface("pixel", cPixel).Msg("on line")
			line = append(line, cPixel)
			continue
		}

		idx := c.FindLine(lines, cPixel, &weight)
		if idx != -1 {
			cl.Debug().Interface("pixel", cPixel).Msg("on another line")
			cl.Debug().Int("onLine", idx).Msg("line pixel was found to be on")
			lines[idx] = append(lines[idx], cPixel)
			continue
		}

		cl.Debug().Interface("pixel", cPixel).Msg("couldn't find a line")
		lines = append(lines, line)

		var newline Line
		line = newline
		line = append(line, cPixel)
	}
	lines = append(lines, line)

	cl.Debug().Msg("Finished")
	return lines
}

// GetNearestBlack returns the black pixel on the running asbuilt that is
// closest to the given x,y poisition
func (c *Canvas) GetNearestBlack(x, y int) (int, int) {
	//cl := c.log.With().Str("func", "DrawLines").Logger()
	bnds := c.img.Bounds()

	minX := bnds.Min.X
	maxX := bnds.Max.X
	minY := bnds.Min.Y
	maxY := bnds.Max.Y

	for i := 0; i <= 15; i++ {
		if x-i < minX || x+i > maxX || y-i < minY || y+i > maxY {
			return x, y
		}
		aboveX := x - i
		aboveY := y + i

		belowX := x + i
		belowY := y - i

		clr := c.GetColor(aboveX, belowY)
		if c.CompareColors(clr, Black) {
			return aboveX, belowY
		}

		clr = c.GetColor(aboveX, aboveY)
		if c.CompareColors(clr, Black) {
			return aboveX, aboveY
		}

		clr = c.GetColor(belowX, aboveY)
		if c.CompareColors(clr, Black) {
			return belowX, aboveY
		}

		clr = c.GetColor(belowX, belowY)
		if c.CompareColors(clr, Black) {
			return belowX, belowY
		}

	}
	return x,y
}

// CompareColors compares two colors to see if they're the same
func (c *Canvas) CompareColors(c1, c2 color.RGBA) bool {
	return (c1.R == c2.R && c1.B == c2.B) && c1.G == c2.G
}

// IntAbs returns the absolute value of an integer
func (c *Canvas) IntAbs(x int) int {
	if x < 0 {
		return x*-1
	}
	return x
}

// OnLine checks to see if a given point is on the current line
func (c *Canvas) OnLine(line Line, pixel *Pixel, weight *float64) bool {
	if len(line) == 0 || len(line) == 1 {
		return true
	}

	x1, y1 := line[len(line)-1].X, line[len(line)-1].Y
	x2, y2 := pixel.X, pixel.Y
	difX := c.IntAbs(x2-x1)
	difY := c.IntAbs(y2-y1)
	if difX == difY {
		return true
	}
	return false
}

// FindLine finds which line the pixel is on
func (c *Canvas) FindLine(lines Lines, pixel *Pixel, weight *float64) int {
	for i, line := range lines {
		if c.OnLine(line, pixel, weight) {
			return i
		}
	}
	return -1
}

// ComparePixels compares two pixels to see if they're on the same line
//
// Return values -
//    -1 : Pixel is *very* far away from previous one
//     1 : Pixel is more than the standard distance from the last one,
//          it's probably on the previous line
//     0 : Pixel is on the same line as the current one
func (c *Canvas) ComparePixels(curr, prev *Pixel) int {
	currFx, prevFx := float64(curr.X), float64(prev.X)
	currFy, prevFy := float64(curr.Y), float64(prev.Y)

	deltaX := math.Abs(currFx - prevFx)
	deltaY := math.Abs(currFy - prevFy)
	if deltaX > xPlus || deltaX < xMinus || deltaY > yPlus || deltaY < yMinus {
		return -1
	}

	if deltaY > 10 {
		return 1
	}
	return 0
}

// DrawAntialiased draws an antialiaed line using Xiaolin Wu's line algorithm
// https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm
func (c *Canvas) DrawAntialiased(start, end Pixel, clr color.RGBA) {
	cl := c.log.With().Str("func", "DrawAntialiased").Logger()
	cl.Debug().Interface("startPixel", start).Msg("Starting pixel")
	cl.Debug().Interface("endPixel", end).Msg("Ending pixel")

	// the math package uses float64s for its parameters, so let's get
	// our x and y values as float64
	x1 := float64(start.X)
	y1 := float64(start.Y)
	x2 := float64(end.X)
	y2 := float64(end.Y)

	// Check if the line is more vertical than horizontal
	steep := math.Abs(y2-y1) > math.Abs(x2-x1)
	if steep {
		// Switch quadrants
		x1, x2 = c.SwapValues(x1, x2)
		y1, y2 = c.SwapValues(y1, y2)
	}

	// Check if we need to switch quadrants based on the horizontal direction
	if x1 > x2 {
		// Switch quadrants
		x1, x2 = c.SwapValues(x1, x2)
		y1, y2 = c.SwapValues(y1, y2)
	}

	// Calculate the slope of the line
	dx := x2 - x1
	dy := y2 - y1
	gradient := 1.0
	if dx != 0 {
		gradient = dy / dx
	}

	// First point on the line
	xend := math.Floor(x1 + 0.5)
	yend := y1 + gradient*(xend-x1)

	intersect := yend + gradient

	// Use the inverse mantissa to calculate the gap
	xgap := 1 - c.Mantissa(x1+0.5)

	xpixel1 := xend
	ypixel1 := math.Floor(yend)

	if steep {
		c.Plot(ypixel1, xpixel1, ypixel1+1, xpixel1, clr, c.Mantissa(yend)*xgap)
	} else {
		c.Plot(xpixel1, ypixel1, xpixel1, ypixel1+1, clr, c.Mantissa(yend)*xgap)
	}

	// Second point on the line
	xend = math.Floor(x2 + 0.5)
	yend = y2 + gradient*(xend-x2)
	xgap = c.Mantissa(x2 + 0.5)
	xpixel2 := xend
	ypixel2 := math.Floor(yend)
	if steep {
		c.Plot(ypixel2, xpixel2, ypixel2+1, xpixel2, clr, c.Mantissa(yend)*xgap)
	} else {
		c.Plot(xpixel2, ypixel2, xpixel2, ypixel2+1, clr, c.Mantissa(yend)*xgap)
	}

	// Now for the actual line itself
	if steep {
		for x := xpixel1 + 1; x <= xpixel2-1; x++ {
			c.Plot(math.Floor(intersect), x, math.Floor(intersect)+1, x, clr, c.Mantissa(intersect))
			intersect = intersect + gradient
		}
	} else {
		for x := xpixel1 + 1; x <= xpixel2-1; x++ {
			c.Plot(x, math.Floor(intersect), x, math.Floor(intersect)+1, clr, c.Mantissa(intersect))
			intersect = intersect + gradient
		}
	}
}

// Plot sets the color of the pixels at the given points
func (c *Canvas) Plot(x1, y1, x2, y2 float64, c1 color.RGBA, weight float64) {
	X1, Y1, X2, Y2 := int(x1), int(y1), int(x2), int(y2)
	weight = math.Min(1.0, weight*1.3)
	c2 := c.GetWeightedColor(c.GetColor(X1, Y1), c1, weight*1.2)

	if canSet, ok := c.img.(CanSet); ok {
		canSet.Set(X1, Y1, c1)
		canSet.Set(X2, Y2, c2)
	}
}

// GetWeightedColor gets the the new color using the given weight value
func (c *Canvas) GetWeightedColor(c1, c2 color.RGBA, weight float64) color.RGBA {
	weight = math.Min(1, math.Max(0, weight))
	color1R, color1G, color1B := float64(c1.R), float64(c1.G), float64(c1.B)
	weightR, weightG, weightB := float64(c2.R)-color1R, float64(c2.G)-color1G, float64(c2.B)-color1B
	return color.RGBA{
		R: uint8(color1R + (weightR * weight)),
		G: uint8(color1G + (weightG * weight)),
		B: uint8(color1B + (weightB * weight)),
		A: c2.A,
	}
}

// Mantissa returns the fraction portion of the float
func (c *Canvas) Mantissa(x float64) float64 {
	return x - math.Floor(x)
}

// SwapValues swaps two floating point values
func (c *Canvas) SwapValues(v1, v2 float64) (float64, float64) {
	return v2, v1
}

/*
// saveApproxChanges saves the approxChanges pixels array to a CSV file
// ... this was used to help me figure out what's going on :D
func (c *Canvas) saveApproxChanges(pixels []Pixel) {
    file, err := os.Create("approxChanges.csv")
    if err != nil {
        fmt.Printf("failed to create csv file\n")
        return
    }
    defer file.Close()

    wr := csv.NewWriter(file)
    defer wr.Flush()

    headers := []string{"X", "Y"}
    wr.Write(headers)

    for _, pixel := range pixels {
        var record []string
        record = append(record, fmt.Sprint(pixel.X))
        record = append(record, fmt.Sprint(pixel.Y))
        wr.Write(record)
    }
}*/
