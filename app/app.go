// Package app provides user input hanlding
package app

import (
	"caddae/imageproc"
	"caddae/types"

	"github.com/jroimartin/gocui"
)

// SetUserInput sets the user input configuration
func (a *App) SetUserInput(input UserInput) {
	al := a.Log.With().Str("func", "SetUserInput").Logger()
	al.Debug().Interface("userinput", input).Send()
	a.in = input
}

// Start starts the application process of validating user input and
// processing the given images.
func (a *App) Start(u types.UI, g *gocui.Gui) error {
	al := a.Log.With().Str("func", "Start").Logger()

	al.Debug().Msg("Checking input values")

	// Validate user input
	conf, err := a.ValidateInput()
	if err != nil {
		return err
	}

	// Let the user know the input was good
	msg := "Input successfully validated!\n"
	g.Update(func(*gocui.Gui) error {
		if err := u.Log(msg); err != nil {
			return err
		}
		return nil
	})
	al.Debug().Msg("valid input")
	al.Debug().Msg("starting image pre processing")

	// Update the user on what we're doing
	msg = "Starting image pre processing..\n"
	g.Update(func(*gocui.Gui) error {
		if err := u.Log(msg); err != nil {
			return err
		}
		return nil
	})

	// Create a new image processor with the given configuration
	a.Ip = imageproc.New(conf, &a.Log)

	// Start image processing
	if err := a.Ip.ProcessImages(u, g); err != nil {
		return err
	}

	return nil
}
