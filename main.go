package main

import (
	"embed"
	"github.com/downace/print-server/internal/gui"
	"github.com/downace/print-server/internal/logging"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "Print Server"

func main() {

	logging.InitLogs()

	err := gui.RunApp(AppName, assets)

	if err != nil {
		println("Error:", err.Error())
	}
}
