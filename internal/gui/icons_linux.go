package gui

import (
	_ "embed"
)

//go:embed resources/trayicon_default.png
var trayIconDefault []byte

//go:embed resources/trayicon_running.png
var trayIconRunning []byte

//go:embed resources/trayicon_stopped.png
var trayIconStopped []byte

//go:embed resources/trayicon_error.png
var trayIconError []byte
