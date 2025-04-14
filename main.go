package main

import (
	"embed"
	"github.com/downace/print-server/internal/logging"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

const AppName = "Print Server"

func main() {
	logging.InitLogs()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:         AppName,
		Width:         400,
		Height:        600,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.startup,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon: appIcon,
		},
		Logger: logging.WailsLog,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
