package app

import (
	"caddae/imageproc"
	"errors"

	"github.com/rs/zerolog"
)

// App object to hold the user input and image processor
type App struct {
	Log zerolog.Logger
	Ip  *imageproc.ImageProc
	in  UserInput
}

// UserInput object to hold input from the UI
type UserInput struct {
	Rl       string `json:"redline"`
	Ra       string `json:"running"`
	Jn       string `json:"job_number"`
	Wpd      string `json:"wpd"`
	Strand   string `json:"strand"`
	Cable    string `json:"cable"`
	Overlash string `json:"overlash"`
	Anchors  string `json:"anchors"`
}

// InvFileErr is the error we'll throw if the user gave us an invalid file type
var InvFileErr error = errors.New("provided file type is not allowed. allowed types are .JSON")
