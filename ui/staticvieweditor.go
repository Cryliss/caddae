package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// Edit is editor function for the static view editor.
func (e *staticViewEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	_, y := v.Cursor()
	maxY := strings.Count(v.Buffer(), "\n")
	switch key {
	case gocui.KeyArrowDown:
		if y < maxY {
			v.MoveCursor(0, 1, true)
		}
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
