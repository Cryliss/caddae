package ui

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

// createPanelView creates the panel view.
func (u *UI) createPanelView(name string, x1, y1, x2, y2 int) (*gocui.View, error) {
	v, err := u.g.SetView(name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return nil, err
	}

	p := panelViews[name]
	v.Title = p.title
	v.Editable = p.edit

	if err := u.write(name, p.body); err != nil {
		return nil, err
	}

	switch name {
	case ABOUT_PANEL:
		fallthrough
	case CREATE_BUTTON:
		fallthrough
	case LOG_PANEL:
		fallthrough
	case PRODUCTION_PANEL:
		v.Highlight = false
		v.Autoscroll = false
		v.Wrap = true
		v.Editor = newEditor(u, &staticViewEditor{})
		break
	case C300_01_PANEL:
		fallthrough
	case C300_02_PANEL:
		fallthrough
	case C300_03_PANEL:
		v.Highlight = true
		v.Autoscroll = false
		v.Editor = newEditor(u, nil)
		v.SelFgColor = gocui.ColorGreen
		break
	case C300_04_PANEL:
		v.Highlight = true
		v.Autoscroll = false
		v.Editor = newEditor(u, nil)
		v.SelFgColor = gocui.ColorMagenta
		break
	default:
		v.Highlight = true
		v.Autoscroll = false
		v.Editor = newEditor(u, nil)
		break
	}
	return v, nil
}

// aactivatePanelView ctivates the view defined by id.
func (u *UI) activatePanelView(id int) error {
	if err := u.setPanelView(u.views[id]); err != nil {
		return err
	}
	v := panelViews[u.views[id]]
	switch v.title {
	case C300_01_PANEL:
		fallthrough
	case C300_02_PANEL:
		fallthrough
	case C300_03_PANEL:
		u.g.SelFgColor = gocui.ColorGreen
		break
	case C300_04_PANEL:
		u.g.SelFgColor = gocui.ColorYellow
		break
	default:
		u.g.SelFgColor = gocui.ColorCyan
	}
	u.g.Cursor = v.cursor
	u.cv = id

	return nil
}

// setPanelView sets the panel view
func (u *UI) setPanelView(name string) error {
	if err := u.closeModal(u.cm); err != nil {
		return err
	}

	// Save cursor position before switch view
	view := u.g.CurrentView()
	x, y := view.Cursor()
	u.c.Set(view.Name(), x, y)

	if _, err := u.g.SetCurrentView(name); err != nil {
		if err == gocui.ErrUnknownView {
			return nil
		}
		return err
	}
	return nil
}

// Write writes the content into the specific view and set the cursor to the buffer end
func (u *UI) write(name, text string) error {
	v, err := u.g.View(name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(v, text)
	v.SetCursor(len(text), 0)
	u.c.Set(name, len(text), 0)

	return nil
}

// readEditView reads the user entered content from an editable view
func (u *UI) readEditView(name string) (string, error) {
	// Make sure the panel that's being requested to read from is actually
	// editable first
	p := panelViews[name]
	if !p.edit {
		e := fmt.Sprintf("u.readView(%s): No content to retrieve; view is not editable!", name)
		return "", errors.New(e)
	}

	// Get the view by name
	v, err := u.g.View(name)
	if err != nil {
		return "", err
	}

	// If we actually got a view, return the views buffer.
	if v != nil {
		content := strings.Split(v.Buffer(), "\n")
		return content[0], nil
	}

	return "", nil
}

// findView finds the view defined by name and returns the view index
func (u *UI) findView(name string) int {
	var viewId = -1
	for idx, v := range u.views {
		if v == name {
			viewId = idx
			break
		}
	}
	return viewId
}

// updateView updates the view content
func (u *UI) updateView(v *gocui.View, buffer string) error {
	if v != nil {
		v.Clear()
		if err := u.write(v.Name(), buffer); err != nil {
			return err
		}
	}
	return nil
}

// nextView activates the next panel
func (u *UI) nextView(wrap bool) error {
	var index int
	index = u.cv + 1
	if index > len(panelViews)-1 {
		if wrap {
			index = 0
		} else {
			return nil
		}
	}
	u.cv = index % len(panelViews)
	return u.activatePanelView(u.cv)
}

// prevView activates the previous panel
func (u *UI) prevView(wrap bool) error {
	var index int
	index = u.cv - 1
	if index < 0 {
		if wrap {
			index = len(panelViews) - 1
		} else {
			return nil
		}
	}
	u.cv = index % len(panelViews)
	return u.activatePanelView(u.cv)
}

// ClearView clears the panel view
func (u *UI) ClearView(name string) {
	v, _ := u.g.View(name)
	v.Clear()
}

// DeleteView deletes the current view
func (u *UI) DeleteView(name string) {
	v, _ := u.g.View(name)
	u.g.DeleteView(v.Name())
}
