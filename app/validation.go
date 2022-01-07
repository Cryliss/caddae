package app

import (
	"caddae/imageproc"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ValidateInput validates the users input values
func (a *App) ValidateInput() (imageproc.Config, error) {
	al := a.Log.With().Str("func", "ValidateInput").Logger()

	var conf imageproc.Config

	// Check redline file path actually exists
	if _, err := os.Stat(a.in.Rl); errors.Is(err, os.ErrNotExist) {
		al.Err(err).Str("redline", a.in.Rl).Msg("ValidateInput")
		e := fmt.Sprintf("a.ValidateInput: provided redline file path '%s' does not exist!", a.in.Rl)
		return conf, errors.New(e)
	}

	// Check the file extension
	fp := strings.Split(a.in.Rl, ".")
	ext := fp[1]
	if ext != "png" {
		e := fmt.Sprintf("a.ValidateInput: incorrect redline file type given! Only .png is allowed, you provided '.%s'", ext)
		return conf, errors.New(e)
	}

	// No issues with the input. Set the confguration value
	conf.Rl = a.in.Rl

	// Check running asbuilt file path actually exists
	if _, err := os.Stat(a.in.Ra); errors.Is(err, os.ErrNotExist) {
		e := fmt.Sprintf("a.ValidateInput: provided running file path '%s' does not exist!", a.in.Ra)
		return conf, errors.New(e)
	}

	// Check the file extension
	fp = strings.Split(a.in.Ra, ".")
	ext = fp[1]
	if ext != "png" {
		e := fmt.Sprintf("a.ValidateInput: incorrect running file type given! Only .png is allowed, you provided '.%s'", ext)
		return conf, errors.New(e)
	}

	// No issues with the input. Set the confguration value
	conf.Ra = a.in.Ra

	// Check thes job number
	jn := strings.Split(a.in.Jn, "_")
	if len(jn) < 2 {
		e := fmt.Sprintf("a.ValidateInput: invalid job number given! '%s'", a.in.Jn)
		return conf, errors.New(e)
	}
	if jn[0] != "DYEA" && jn[1] != "LSA" {
		if jn[0] != "VZ" && jn[1] != "LAN" {
			e := fmt.Sprintf("a.ValidateInput: invalid job number given! '%s'", a.in.Jn)
			return conf, errors.New(e)
		}
	}

	// No issues with the input. Set the confguration value
	conf.Jn = a.in.Jn

	// Check the work performed date
	_, err := time.Parse("01/02/2006", a.in.Wpd)
	if err != nil {
		e := fmt.Sprintf("a.ValidateInput: error parsing workdate!\ntime.Parse('01/02/2006', '%s'): %v", a.in.Wpd, err)
		return conf, errors.New(e)
	}

	// No issues with the input. Set the confguration value
	conf.Wpd = a.in.Wpd

	// Now let's check the unit values
	if a.in.Strand != "" {
		strand, err := strconv.ParseFloat(a.in.Strand, 64)
		if err != nil {
			e := fmt.Sprintf("a.ValidateInput: invalid quantity given for C300-01! '%s'", a.in.Strand)
			return conf, errors.New(e)
		}
		conf.Strand = strand
	}

	if a.in.Cable != "" {
		cable, err := strconv.ParseFloat(a.in.Cable, 64)
		if err != nil {
			e := fmt.Sprintf("a.ValidateInput: invalid quantity given for C300-02! '%s'", a.in.Cable)
			return conf, errors.New(e)
		}
		conf.Cable = cable
	}

	if a.in.Overlash != "" {
		overlash, err := strconv.ParseFloat(a.in.Overlash, 64)
		if err != nil {
			e := fmt.Sprintf("a.ValidateInput: invalid quantity given for C300-03! '%s'", a.in.Overlash)
			return conf, errors.New(e)
		}
		conf.Overlash = overlash
	}

	if a.in.Anchors != "" {
		anchors, err := strconv.ParseFloat(a.in.Anchors, 64)
		if err != nil {
			e := fmt.Sprintf("a.ValidateInput: invalid quantity given for C300-04! '%s'", a.in.Anchors)
			return conf, errors.New(e)
		}
		conf.Anchors = anchors
	}

	return conf, nil
}
