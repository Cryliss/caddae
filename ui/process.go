package ui

import (
	"caddae/app"
	"fmt"
	"image"
	"image/draw"
	"strings"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
	"github.com/jroimartin/gocui"
)

// checkUserInput checks the provided user input.
func (u *UI) checkUserInput() error {
	// Get the user entered values
	redline, err := u.readEditView(REDLINE_PANEL)
	if err != nil {
		return err
	}

	running, err := u.readEditView(RUNNING_PANEL)
	if err != nil {
		return err
	}

	job, err := u.readEditView(JOB_PANEL)
	if err != nil {
		return err
	}
	job = strings.ToUpper(job)

	wpd, err := u.readEditView(WPD_PANEL)
	if err != nil {
		return err
	}

	strand, err := u.readEditView(C300_01_PANEL)
	if err != nil {
		return err
	}

	cable, err := u.readEditView(C300_02_PANEL)
	if err != nil {
		return err
	}

	overlash, err := u.readEditView(C300_03_PANEL)
	if err != nil {
		return err
	}

	anchors, err := u.readEditView(C300_04_PANEL)
	if err != nil {
		return err
	}

	in := app.UserInput{
		Rl:       redline,
		Ra:       running,
		Jn:       job,
		Wpd:      wpd,
		Strand:   strand,
		Cable:    cable,
		Overlash: overlash,
		Anchors:  anchors,
	}

	u.a.SetUserInput(in)
	return nil
}

// createRunning is triggered by pressing Create Running AsBuilt button
// Validates user input, pre processes the image and creates the running asbuilt.
func (u *UI) createRunning(g *gocui.Gui, v *gocui.View) error {
	u.mu.Lock()
	if u.started {
		u.mu.Unlock()
		return nil
	}
	u.started = true
	u.mu.Unlock()

	if err := u.checkUserInput(); err != nil {
		u.LogErr(fmt.Sprintf("%v", err))
		u.mu.Lock()
		u.started = false
		u.mu.Unlock()
		return nil
	}

	u.ClearLog()

	go func() {
		err := u.a.Start(u, g)
		defer func() {
			u.mu.Lock()
			u.started = false
			u.mu.Unlock()
			if err != nil {
				msg := fmt.Sprintf("%v", err)
				u.LogErr(msg)
			}
			if err == nil {
				// TODO: Make this work ???
				gl.StartDriver(func(driver gxui.Driver) {
					file := u.a.Ip.Running()
					theme := flags.CreateTheme(driver)
					img := theme.CreateImage()

					dx, dy := file.Bounds().Max.X, file.Bounds().Max.Y
					if dx < MinWindowSize {
						dx = MinWindowSize
					}
					if dy < MinWindowSize {
						dy = MinWindowSize
					}

					window := theme.CreateWindow(dx, dy, "Running AsBuilt Preview")
					window.SetScale(flags.DefaultScaleFactor)
					window.AddChild(img)

					rgba := image.NewRGBA(file.Bounds())
					draw.Draw(rgba, file.Bounds(), file, image.ZP, draw.Src)
					texture := driver.CreateTexture(rgba, 1)
					img.SetTexture(texture)

					window.OnClose(driver.Terminate)
				})
			}
		}()
	}()

	return nil
}
