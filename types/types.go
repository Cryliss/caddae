package types

import "image/color"

// Package to help us avoid import cycles in Go.

// UI is the user nterface for the UI so we can perform logging
type UI interface {
	Log(message string) error
	LogErr(message string) error
	ClearLog() error
}

// Production details for the callout box
type Production struct {
	Date  string
	Units []Unit
}

// Unit details for the production box
type Unit struct {
	Name  string
	Qty   string
	Text  string
	Color color.RGBA
}
