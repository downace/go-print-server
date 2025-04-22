package main

import (
	"embed"
	"github.com/downace/print-server/internal/cli"
	"github.com/downace/print-server/internal/gui"
	"github.com/downace/print-server/internal/logging"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "Print Server"

func main() {

	logging.InitLogs()

	isCli := len(os.Args) >= 2 && os.Args[1] == "cli"

	var err error
	if isCli {
		err = cli.RunApp(os.Args[2:])
	} else {
		// https://github.com/wailsapp/wails/issues/2977
		_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
		err = gui.RunApp(AppName, assets)
	}

	if err != nil {
		println("Error:", err.Error())
	}
}
