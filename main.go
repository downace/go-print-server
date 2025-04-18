package main

import (
	"embed"
	"github.com/downace/print-server/internal/cli"
	"github.com/downace/print-server/internal/gui"
	"github.com/downace/print-server/internal/logging"
	"os"
	"slices"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "Print Server"

func main() {

	logging.InitLogs()

	isCli := slices.Index(os.Args, "--cli") >= 0

	var err error
	if isCli {
		err = cli.RunApp()
	} else {
		// https://github.com/wailsapp/wails/issues/2977
		_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
		err = gui.RunApp(AppName, assets)
	}

	if err != nil {
		println("Error:", err.Error())
	}
}
