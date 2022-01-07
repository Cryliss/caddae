package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

// OpenModal creates and opens the modal window. If "autoHide" parameter is true,
// the modal will be automatically closed after 2 minutes.
func (u *UI) OpenModal(name string, w, h int, autoHide bool) (*gocui.View, error) {
	_, ok := modalViews[name]
	if !ok {
		e := fmt.Sprintf("u.OpenModal(%s, %d, %d, %v): Modal not found.", name, w, h, autoHide)
		return nil, errors.New(e)
	}

	v, err := u.createModal(name, 20, 1, w, h)
	if err != nil {
		return nil, err
	}

	if err := u.setPanelView(name); err != nil {
		return nil, err
	}
	u.cm = name

	if autoHide && u.cm == HELP_PANEL {
		// Close the modal automatically after 5 seconds
		u.dt = time.AfterFunc(5*time.Second, func() {
			u.g.Update(func(*gocui.Gui) error {
				if err := u.closeModal(name); err != nil {
					return err
				}
				return nil
			})
		})
	} else if autoHide {
		// Close the modal automatically after 2 minutes
		u.dt = time.AfterFunc(2*time.Minute, func() {
			u.g.Update(func(*gocui.Gui) error {
				if err := u.closeModal(name); err != nil {
					return err
				}
				return nil
			})
		})
	}
	return v, nil
}

// closeOpenedModals closes all the opened modal elements.
func (u *UI) closeOpenedModals(modals []string) error {
	for _, m := range modals {
		if view, _ := u.g.View(m); view != nil {
			u.closeModal(view.Name())
		}
	}
	return nil
}

// closeModal closes the modal window and restores the focus to the last accessed panel view.
func (u *UI) closeModal(modals ...string) error {
	for _, name := range modals {
		if _, err := u.g.View(name); err != nil {
			if err == gocui.ErrUnknownView {
				return nil
			}
			return err
		}
		u.g.DeleteView(name)
		u.g.DeleteKeybindings(name)
		u.g.Cursor = true
		u.cm = ""
	}
	return u.activatePanelView(u.cv)
}

// createModal creates the modal view.
func (u *UI) createModal(name string, x1, y1, x2, y2 int) (*gocui.View, error) {
	v, err := u.g.SetView(name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return nil, err
	}
	m := modalViews[name]

	v.Title = m.title
	v.Editable = m.edit

	if err := u.write(name, m.body); err != nil {
		return nil, err
	}
	return v, nil
}
