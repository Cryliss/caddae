package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// MakeCursors creates a new map of cursors and set their initial x,y values
func MakeCursors(numViews int) Cursors {
	// Make a new map of cursors, with the same length as the views
	return make(Cursors, numViews)
}

// Restore restores the cursor position of the provided view
func (c Cursors) Restore(view *gocui.View) error {
	return view.SetCursor(c.Get(view.Name()))
}

// Get returns the current position of the view
func (c Cursors) Get(view string) (int, int) {
	if v, ok := c[view]; ok {
		return v.x, v.y
	}
	return 0, 0
}

// Set the x and y positions of the current views cursor
func (c Cursors) Set(view string, x, y int) {
	if v, ok := c[view]; ok {
		v.x = x
		v.y = y
	}
}

// scrollDown moves the cursor to the next buffer line.
func (u *UI) scrollDown(g *gocui.Gui, v *gocui.View) error {
	maxY := strings.Count(v.Buffer(), "\n")
	if maxY < 1 {
		v.SetCursor(0, 0)
	}
	return nil
}
