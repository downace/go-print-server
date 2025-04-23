package gui

import (
	_ "embed"
)

//go:embed resources/trayicon_default.ico
var trayIconDefault []byte

//go:embed resources/trayicon_running.ico
var trayIconRunning []byte

//go:embed resources/trayicon_stopped.ico
var trayIconStopped []byte

//go:embed resources/trayicon_error.ico
var trayIconError []byte
