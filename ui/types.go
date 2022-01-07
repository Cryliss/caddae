package ui

import (
	"caddae/app"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/rs/zerolog"
)

// BEGIN bindings.go Types {{{

// ClosureFn is the function we'll use as a closure function
type ClosureFn func(*gocui.Gui, *gocui.View) error

// handler
type handler struct {
	views   []string
	key     interface{}
	keyName string
	help    string
	action  func(*UI, bool) ClosureFn
}

var handlers = []handler{
	{views, gocui.KeyTab, "Tab", "Next Panel", nextPanel},
	{nil, gocui.KeyPgup, "PgUp", "Jump to the top", nil},
	{nil, gocui.KeyPgdn, "PgDown", "Jump to the bottom", nil},
	{nil, gocui.KeyHome, "Home", "Jump to the start", nil},
	{nil, gocui.KeyEnd, "End", "Jump to the end", nil},
	{nil, gocui.KeyCtrlC, "Ctrl+c", "Quit", quit},
	{nil, gocui.KeyCtrlX, "Ctrl+x", "Clear editor content", nil},
	{nil, gocui.KeyCtrlZ, "Ctrl+z", "Restore editor content", nil},
}

// END bindings.go Types }}}

// BEGIN cursors.go Types {{{

// Cursors is a map to store cursor positions in each view so we can restore it
// once the view is active again
type Cursors map[string]struct{ x, y int }

// END cursors.go Types }}}

// BEGIN editor.go Types {{{

// editor is a ui editor
type editor struct {
	ui     *UI
	editor gocui.Editor
}

// END editor.go Types }}}

// BEGIN modal.go Types {{{

const helpText = `Instructions for Use:
=====================
Please update each of the following widgets with
the requested information

Redline
-------
Enter the full path name of the redline file you wish to
digitally recreate.

Format: .png

Running AsBuilt
---------------
Enter the full path name of the original asbuilt file
that is to be updated.

Format: .png

DYEA/VZ
-------
Enter the DYEA/VZ# associated with the provided redline.

Format: .DYEA_LSA_8XXXXXX | VZ_LAN_0000XXXX

WPD
---
Enter the date the work was performed.

Format: MM/DD/YYYY

Production
-----------
Enter the each applicable production units quantities.

Format: 100 | 85.25

Keybindings
===========
`

var modals = []string{
	HELP_PANEL,
	DIAGRAM_MODAL,
	PROGRESS_MODAL,
}

var modalViews = map[string]Panel{
	HELP_PANEL: {
		title: "CADDAE Help",
		body:  "",
		edit:  false,
	},
	DIAGRAM_MODAL: {
		title:  "Resulting Running AsBuilt",
		body:   "",
		edit:   false,
		cursor: false,
	},
}

// END modal.go Types }}}

// staticViewEditor is an editor  for static (non-editable) views
type staticViewEditor editor

// END staticViewEditor.go Types }}}

// BEGIN ui.go Types {{{

// Ui struct
type UI struct {
	a       *app.App
	g       *gocui.Gui
	l       zerolog.Logger
	mu      sync.Mutex
	cv      int         // The currently active panel
	nv      int         // The next panel in the array
	cm      string      // The currently active modal
	c       Cursors     // Tracking for cursor positions in each views
	dt      *time.Timer // Timer for displaying the diagram
	lt      *time.Timer // Timer for logging the recreation process
	lm      []string    // Array of log messages
	started bool        // Whether or not the process has been started.
	views   []string    // Array of views
}

// END ui.go Types }}}

// BEGIN views.go Types {{{

const aboutText = `This application is meant to read a *simple* aerial redline file and recreate it digitally as a running asbuilt.

Please fill out the boxes below and hit 'Create Running AsBuilt! when you're ready to begin.
After that, you can view what's happening during the process in the 'Log' panel.

Upon completion, it will save the running asbuilt in the 'edits' folder for your review.

Press Ctrl+H to toggle the help modal.
`

const prodText = `Provide quantities for the applicable units















`

// Panel object for the UI
type Panel struct {
	title  string
	body   string
	x1, y1 int
	x2, y2 int
	cache  []byte
	edit   bool
	cursor bool
	editor *UI
}

// MinWindowSize for the terminal
const MinWindowSize = 125

// Panel name constants.
const (
	ABOUT_PANEL      = "about"
	HELP_PANEL       = "help"
	REDLINE_PANEL    = "redline"
	RUNNING_PANEL    = "running"
	JOB_PANEL        = "job"
	WPD_PANEL        = "wpd"
	PRODUCTION_PANEL = "prod"
	C300_01_PANEL    = "c300_01"
	C300_02_PANEL    = "c300_02"
	C300_03_PANEL    = "c300_03"
	C300_04_PANEL    = "c300_04"
	DIAGRAM_MODAL    = "diagram"
	PROGRESS_MODAL   = "progress"
	LOG_PANEL        = "log"
	CREATE_BUTTON    = "create"
)

// Initialize the array of view names.
var views = []string{
	ABOUT_PANEL,
	REDLINE_PANEL,
	RUNNING_PANEL,
	JOB_PANEL,
	WPD_PANEL,
	PRODUCTION_PANEL,
	C300_01_PANEL,
	C300_02_PANEL,
	C300_03_PANEL,
	C300_04_PANEL,
	CREATE_BUTTON,
	LOG_PANEL,
}

// Initialize the panel views for the UI and save them in a map.
var panelViews = map[string]Panel{
	ABOUT_PANEL: {
		title:  "CADDAE",
		body:   aboutText,
		x1:     0,
		y1:     0,
		x2:     120,
		y2:     9,
		edit:   true,
		cursor: true,
	},
	REDLINE_PANEL: {
		title:  "Redline",
		body:   "",
		x1:     0,
		y1:     10,
		x2:     120,
		y2:     12,
		edit:   true,
		cursor: true,
	},
	RUNNING_PANEL: {
		title:  "Running AsBuilt",
		body:   "",
		x1:     0,
		y1:     13,
		x2:     120,
		y2:     15,
		edit:   true,
		cursor: true,
	},
	JOB_PANEL: {
		title:  "DYEA/VZ#",
		body:   "",
		x1:     51,
		y1:     16,
		x2:     120,
		y2:     18,
		edit:   true,
		cursor: true,
	},
	WPD_PANEL: {
		title:  "WPD",
		body:   "",
		x1:     51,
		y1:     19,
		x2:     120,
		y2:     21,
		edit:   true,
		cursor: true,
	},
	PRODUCTION_PANEL: {
		title:  "Production",
		body:   prodText,
		x1:     0,
		y1:     16,
		x2:     50,
		y2:     26,
		edit:   false,
		cursor: true,
	},
	C300_01_PANEL: {
		title:  "C300-01",
		body:   "",
		x1:     3,
		y1:     19,
		x2:     23,
		y2:     21,
		edit:   true,
		cursor: true,
	},
	C300_02_PANEL: {
		title:  "C300-02",
		body:   "",
		x1:     25,
		y1:     19,
		x2:     45,
		y2:     21,
		edit:   true,
		cursor: true,
	},
	C300_03_PANEL: {
		title:  "C300-03",
		body:   "",
		x1:     3,
		y1:     23,
		x2:     23,
		y2:     25,
		edit:   true,
		cursor: true,
	},
	C300_04_PANEL: {
		title:  "C300-04",
		body:   "",
		x1:     25,
		y1:     23,
		x2:     45,
		y2:     25,
		edit:   true,
		cursor: true,
	},
	LOG_PANEL: {
		title:  "Application Log",
		body:   "",
		x1:     0,
		y1:     27,
		x2:     120,
		y2:     45,
		edit:   false,
		cursor: true,
	},
	CREATE_BUTTON: {
		title:  "",
		body:   "                Create running asbuilt! (click me!)",
		x1:     51,
		y1:     24,
		x2:     120,
		y2:     26,
		edit:   false,
		cursor: true,
	},
}

// END views.go Types }}}
