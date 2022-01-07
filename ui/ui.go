package ui

import (
	"caddae/app"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog"
)

// New creates and returns a new Ui
func New(a *app.App, logger *zerolog.Logger) *UI {
	// Initialize a new GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	// Initialize & return a new Ui struct
	ui := UI{
		a:       a,
		g:       g,
		l:       logger.With().Str("module", "ui").Logger(),
		started: false,
		views:   views,
	}

	fl := ui.l.With().Str("func", "New").Logger()
	fl.Debug().Msg("Created")
	return &ui
}

// StartUI initializes the interfaces values, creates keybindings and
// starts the UI's main loop
func (u *UI) StartUI() {
	fl := u.l.With().Str("func", "StartUI").Logger()
	fl.Debug().Msg("Started")

	// Set our default panel settings
	u.g.Highlight = true
	u.g.SelFgColor = gocui.ColorCyan

	// Enable cursor and mouse functionality
	u.g.Cursor = true
	u.g.Mouse = true

	u.c = MakeCursors(len(u.views))
	u.cv = u.findView(REDLINE_PANEL)
	u.nv = 0

	// Set the function to manage the layout of the UI
	u.g.SetManagerFunc(u.Layout)

	// Set the keybindings of the UI
	if err := u.keybindings(); err != nil {
		log.Panicln(err)
	}

	// Start the main loop of the UI
	if err := u.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Layout initialize the panel views and associates the key bindings to them.
func (u *UI) Layout(g *gocui.Gui) error {
	initPanel := func(g *gocui.Gui, v *gocui.View) error {
		// Disable panel views selection with mouse in case the modal is activated
		if u.cm == "" {
			cx, cy := v.Cursor()
			line, err := v.Line(cy)
			if err != nil {
				u.c.Restore(v)
				u.setPanelView(v.Name())
			}

			if cx > len(line) {
				v.SetCursor(u.c.Get(v.Name()))
				u.c.Set(v.Name(), u.getRowCharacterCount(v, cy), cy)
			}
			u.cv = u.findView(v.Name())
			u.setPanelView(v.Name())
			view := panelViews[v.Name()]
			u.g.Cursor = view.cursor
		}
		return nil
	}

	for _, view := range views {
		if view == CREATE_BUTTON {
			if err := u.g.SetKeybinding(view, gocui.KeyEnter, gocui.ModNone, u.createRunning); err != nil {
				return err
			}
			if err := u.g.SetKeybinding(view, gocui.MouseLeft, gocui.ModNone, initPanel); err != nil {
				return err
			}
			if err := u.g.SetKeybinding(view, gocui.MouseRelease, gocui.ModNone, u.createRunning); err != nil {
				return err
			}
		} else {
			if err := u.g.SetKeybinding(view, gocui.MouseLeft, gocui.ModNone, initPanel); err != nil {
				return err
			}
		}

		v := panelViews[view]
		if _, err := u.createPanelView(view, v.x1, v.y1, v.x2, v.y2); err != nil {
			return err
		}
	}

	// Activate the first panel on first run
	if v := u.g.CurrentView(); v == nil {
		_, err := u.g.SetCurrentView(REDLINE_PANEL)
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

// Close closes the gocui.Gui
func (u *UI) Close() {
	u.g.Close()
}
