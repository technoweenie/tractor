package agent

import (
	"github.com/getlantern/systray"
)

type Workspace struct {
	Name       string
	Path       string
	SocketPath string
	MenuItem   *systray.MenuItem
}
