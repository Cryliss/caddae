package drawing

import "image/color"

// Colors we want to use.

// Blue (Cornflower Blue)
var Blue = color.RGBA{44, 149, 237, 255}

// Red is easier to see as a drawn on line whoops
var Red = color.RGBA{242, 55, 4, 255}

// Black color variable
var Black = color.RGBA{0, 0, 0, 255}

// Coral color variable
var Coral = color.RGBA{255, 127, 80, 255}

// White color variable
var White = color.RGBA{255, 255, 255, 255}

// Transparent color variable
var Transparent = color.RGBA{0, 0, 0, 0}

const (
	// BLACKISH is the name of our "blackish" pixel ranges
	BLACKISH = "blackish"
	// WHITEISH is the name of our "whiteish" pixel ranges
	WHITEISH = "whiteish"
	// YELLOWISH is the name of our "yellowish" pixel ranges
	YELLOWISH = "yellowish"
)

// Ranges is the ranges we're going to check during preprocesing
var Ranges = ColorRanges{
	{
		// Because this isn't really "white", its colors
		// we want to identify as white-ish, like greys and such
		Name: WHITEISH,
		// Pure white, #ffffff
		RMax: 0xff, GMax: 0xff, BMax: 0xff,
		// Minimum to fit our "white" is really a grew, #e0e0e0
		RMin: 0xdf, GMin: 0xe3, BMin: 0xe2,
		Replace: false,
		Make:    color.RGBA{0xff, 0xff, 0xff, 0xff},
	},
	{
		Name: YELLOWISH,
		// "pure" (?) yellow, #ffff00
		RMax: 0xfe, GMax: 0xff, BMax: 0xaf,
		// Minimum to fit our "yellow"
		RMin: 0xc0, GMin: 0xb9, BMin: 0x02,
		Replace: false,
		Make:    color.RGBA{0x64, 0x95, 0xed, 0xff}, // change it to blue?
	},
	{
		Name: BLACKISH,
		// Maximum to fit our "black"
		RMax: 0x39, GMax: 0x39, BMax: 0x39,
		// Pure black, #000000
		RMin: 0x00, GMin: 0x00, BMin: 0x00,
		Replace: false,
		Make:    color.RGBA{0x00, 0x00, 0x00, 0xff},
	},
}
