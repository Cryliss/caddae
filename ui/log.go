package ui

// Log adds green color to a log message before writing it to the user.
func (u *UI) Log(message string) error {
	colorMsg := "\x1b[0;32m" + message + "\n"

	u.mu.Lock()
	lm := u.lm
	u.mu.Unlock()

	if len(lm) > 0 && colorMsg == lm[len(lm)-1] {
		return nil
	}

	if len(lm) == 23 {
		lm = append(lm[11:], colorMsg)
	} else {
		lm = append(lm, colorMsg)
	}
	u.mu.Lock()
	u.lm = lm
	u.mu.Unlock()

	var outMsg string
	for _, msg := range lm {
		outMsg = outMsg + msg
	}
	return u.write(LOG_PANEL, outMsg)
}

// LogErr adds red color to a log message before writing it to the user.
func (u *UI) LogErr(message string) error {
	colorMsg := "\x1b[0;31m" + message + "\n"

	u.mu.Lock()
	lm := u.lm
	u.mu.Unlock()

	if len(lm) > 1 && colorMsg == lm[len(lm)-1] {
		return nil
	}

	if len(lm) == 23 {
		lm = append(lm[11:], colorMsg)
	} else {
		lm = append(lm, colorMsg)
	}
	u.mu.Lock()
	u.lm = lm
	u.mu.Unlock()

	var outMsg string
	for _, msg := range lm {
		outMsg = outMsg + msg
	}
	return u.write(LOG_PANEL, outMsg)
}

// ClearLog clears the contents on the log panel.
func (u *UI) ClearLog() error {
	var lm []string
	u.mu.Lock()
	u.lm = lm
	u.mu.Unlock()

	return u.write(LOG_PANEL, "")
}
