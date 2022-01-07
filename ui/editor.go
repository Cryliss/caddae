package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// newEditor ceates a new GUI editor.
func newEditor(ui *UI, handler gocui.Editor) *editor {
	if handler == nil {
		handler = gocui.DefaultEditor
	}
	return &editor{ui, handler}
}

var cache []byte

// Edit is the editor function for editable views.
func (e *editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// Prevent infinite scrolling
	if (key == gocui.KeyArrowDown || key == gocui.KeyArrowRight) && mod == gocui.ModNone {
		_, cy := v.Cursor()
		if _, err := v.Line(cy); err != nil {
			return
		}
	}

	switch key {
	// Disable line wrapping (right arrow key at line end wraps too)
	case gocui.KeyArrowRight:
		cx, cy := v.Cursor()
		// Get the total number of rows in the current view
		maxY := strings.Count(v.Buffer(), "\n")
		// Check if the cursor is on the last row of the current view
		if cy == maxY-1 {
			// Prevent line wrapping on last row
			if cx >= e.ui.getRowCharacterCount(v, cy) {
				return
			}
		}
	case gocui.KeyHome:
		_, cy := v.Cursor()
		v.SetCursor(0, cy)
	case gocui.KeyEnd:
		_, cy := v.Cursor()
		maxX := e.ui.getRowCharacterCount(v, cy)
		v.SetCursor(maxX, cy)
	case gocui.KeyPgup:
		vx, vy := v.Origin()
		if err := v.SetCursor(0, 0); err != nil && vy > 0 {
			if err := v.SetOrigin(vx, 0); err != nil {
				return
			}
		}
	case gocui.KeyPgdn:
		maxX := e.ui.getLastRowCharacterCount(v)
		maxY := strings.Count(v.ViewBuffer(), "\n") - 1
		v.SetCursor(maxX, maxY)
	case gocui.KeyCtrlX:
		if e.ui.isEditableView(v.Name()) {
			cache = []byte(v.ViewBuffer())
			e.ui.setViewCache(v.Name(), cache)
			e.ui.ClearView(v.Name())
			v.SetCursor(0, 0)
		}
	case gocui.KeyCtrlZ:
		if e.ui.isEditableView(v.Name()) {
			cache = e.ui.getViewCache(v.Name())
			if len(cache) > 0 {
				v.Write(cache)
				cache = []byte{}
			}
		}
	}
	e.editor.Edit(v, key, ch, mod)
}

// isEditableView returns whether or not the view is editable.
func (u *UI) isEditableView(name string) bool {
	for key, view := range panelViews {
		if key == name {
			return view.edit
		}
	}
	return false
}

// getViewCache gets the specified views cache.
func (u *UI) getViewCache(name string) []byte {
	for key, view := range panelViews {
		if key == name {
			if view.cache == nil {
				break
			}
			return view.cache
		}
	}
	return []byte{}
}

// setViewCache sets the view cache value.
func (u *UI) setViewCache(name string, cache []byte) {
	for key, view := range panelViews {
		if key == name {
			view.cache = cache
		}
	}
}

// getRowContent returns the row content defined by "y".
func (u *UI) getRowContent(v *gocui.View, y int) []string {
	var row string
	rows := []string{}
	buffer := v.ViewBuffer()

	for _, char := range []byte(buffer) {
		if string(char) == "\n" {
			rows = append(rows, row)
			row = ""
		} else {
			row = row + string(char)
		}
	}
	if len(rows) > 0 && (y > -1 && y < len(rows)) {
		return []string{rows[y]}
	}
	return []string{""}
}

// getLastRowContent returns the last row content.
func (u *UI) getLastRowContent(v *gocui.View) []string {
	var row string
	rows := []string{}
	buffer := v.ViewBuffer()

	for _, char := range []byte(buffer) {
		if string(char) == "\n" {
			rows = append(rows, row)
			row = ""
		} else {
			row = row + string(char)
		}
	}

	if len(rows) > 0 {
		// Traverse up the string slice and remove all the trailing spaces from the end of the text.
		fn := func(rows []string) int {
			var idx = 1
			for {
				current := string(rows[len(rows)-idx:][0])
				if current == "" {
					idx++
				} else {
					break
				}
			}
			return idx
		}
		index := fn(rows)
		return rows[len(rows)-index:]
	}
	return []string{""}
}

// getRowCharacterCount returns the number of characters in the row defined by "y".
func (u *UI) getRowCharacterCount(v *gocui.View, y int) int {
	row := u.getRowContent(v, y)
	return len(strings.Split(row[0], ""))
}

// getLastRowCharacterCount returns the number of characters in the last row.
func (u *UI) getLastRowCharacterCount(v *gocui.View) int {
	lastRow := u.getLastRowContent(v)
	return len(strings.Split(lastRow[0], ""))
}

// getTotalRows returns the total number of rows of the current view.
func (u *UI) getTotalRows(v *gocui.View) int {
	var rows int
	buffer := v.ViewBuffer()

	for _, char := range []byte(buffer) {
		if string(char) == "\n" {
			rows++
		}
	}
	return rows
}

// getPartialViewBuffer returns the view buffer down until the row defined by "n".
func (u *UI) getPartialViewBuffer(v *gocui.View, n int) string {
	var row string
	var idx int
	var newBuffer string

	rows := []string{}
	buffer := v.ViewBuffer()

	for _, char := range []byte(buffer) {
		if string(char) == "\n" {
			rows = append(rows, row)
			row = ""
			if idx > n {
				break
			}
			idx++
		} else {
			row = row + string(char)
		}
	}
	if idx < n {
		newBuffer = strings.Join(rows[:idx], "\n")
	} else {
		newBuffer = strings.Join(rows[:n], "\n")
	}
	return newBuffer
}
