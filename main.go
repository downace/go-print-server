//go:build !cli

package main

import (
	"embed"
	"github.com/downace/print-server/internal/gui"
	"github.com/downace/print-server/internal/logging"
	"os"
	"runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "Print Server"

func main() {
	logging.InitLogs()

	if runtime.GOOS == "linux" {
		// https://github.com/wailsapp/wails/issues/2977
		_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
	}
	err := gui.RunApp(AppName, assets)

	if err != nil {
		println("Error:", err.Error())
	}
}
