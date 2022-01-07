package ui

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jroimartin/gocui"
)

// keybindings sets the ui key bindings
func (u *UI) keybindings() error {
	for _, h := range handlers {
		if len(h.views) == 0 {
			h.views = []string{""}
		}
		if h.action == nil {
			continue
		}
		for _, view := range h.views {
			if err := u.g.SetKeybinding(view, h.key, gocui.ModNone, h.action(u, true)); err != nil {
				return err
			}
		}
	}

	onDown := func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := v.Cursor()

		if cy < u.getTotalRows(v)-1 {
			v.SetCursor(cx, cy+1)
		}
		u.updateView(v, v.Buffer())
		return nil
	}

	onUp := func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := v.Cursor()

		if cy > 0 {
			v.SetCursor(cx, cy-1)
		}
		u.updateView(v, v.Buffer())
		return nil
	}

	if err := u.g.SetKeybinding(LOG_PANEL, gocui.KeyArrowDown, gocui.ModNone, onDown); err != nil {
		return err
	}

	if err := u.g.SetKeybinding(LOG_PANEL, gocui.KeyArrowUp, gocui.ModNone, onUp); err != nil {
		return err
	}

	return u.g.SetKeybinding("", gocui.KeyCtrlH, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return u.toggleHelp(g, HelpContent())
	})
	return nil
}

// nextPanel retrieves the next panel.
func nextPanel(ui *UI, wrap bool) ClosureFn {
	return func(*gocui.Gui, *gocui.View) error {
		return ui.nextView(wrap)
	}
}

// prevPanel retrieves the previous panel
func prevPanel(ui *UI, wrap bool) ClosureFn {
	return func(*gocui.Gui, *gocui.View) error {
		return ui.prevView(wrap)
	}
}

// quit exits the application
func quit(ui *UI, wrap bool) ClosureFn {
	return func(*gocui.Gui, *gocui.View) error {
		return gocui.ErrQuit
	}
}

// toggleHelp toggles the help view on key pressing.
func (u *UI) toggleHelp(g *gocui.Gui, content string) error {
	if err := u.closeOpenedModals(modals); err != nil {
		return err
	}

	panelHeight := strings.Count(content, "\n")

	if u.cm == HELP_PANEL {
		u.g.DeleteKeybinding("", gocui.MouseLeft, gocui.ModNone)
		u.g.DeleteKeybinding("", gocui.MouseRelease, gocui.ModNone)

		// Stop modal timer from firing in case the modal was closed manually.
		// This is needed to prevent the modal being closed before the predefined delay.
		if u.dt != nil {
			u.dt.Stop()
		}
		return u.closeModal(u.cm)
	}

	v, err := u.OpenModal(HELP_PANEL, 80, panelHeight, true)
	if err != nil {
		return err
	}
	u.g.Cursor = false
	v.Editor = newEditor(u, &staticViewEditor{})

	fmt.Fprintf(v, content)
	return nil
}

// HelpContent populates the help panel.
func HelpContent() string {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(w, helpText)
	for _, handler := range handlers {
		if handler.keyName == "" || handler.help == "" {
			continue
		}
		fmt.Fprintf(w, "  %s\t: %s\n", handler.keyName, handler.help)
	}

	fmt.Fprintf(w, "  %s\t: %s\n", "Ctrl+h", "Toggle Help")
	w.Flush()

	return buf.String()
}
