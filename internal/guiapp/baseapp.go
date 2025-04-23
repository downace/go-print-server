package guiapp

import (
	"context"
	"fyne.io/systray"
	"github.com/downace/print-server/internal/common"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	goruntime "runtime"
)

type BaseApp struct {
	Ctx context.Context

	InitTray       func()
	OnWindowHidden func(hidden bool)

	trayStart             func()
	trayEnd               func()
	hidden                bool
	shouldShutdownOnClose bool
}

func (a *BaseApp) Startup(ctx context.Context) {
	a.Ctx = ctx

	if a.InitTray != nil {
		sigintCh := make(chan os.Signal, 1)
		signal.Notify(sigintCh, os.Interrupt)
		common.ListenChannel(sigintCh, func(_ os.Signal) {
			a.shouldShutdownOnClose = true
		})

		a.trayStart, a.trayEnd = systray.RunWithExternalLoop(a.InitTray, nil)
		a.trayStart()
	}
}

func (a *BaseApp) BeforeClose(ctx context.Context) (prevent bool) {
	if a.InitTray == nil {
		return false
	}
	if a.shouldShutdownOnClose {
		return false
	}
	a.SetWindowHidden(true)
	return true
}

func (a *BaseApp) Quit() {
	a.shouldShutdownOnClose = true
	runtime.Quit(a.Ctx)
}

func (a *BaseApp) Shutdown(_ context.Context) {
	if a.trayEnd != nil {
		a.trayEnd()
	}
}

func (a *BaseApp) SetWindowHidden(hidden bool) {
	a.hidden = hidden
	if hidden {
		runtime.WindowHide(a.Ctx)
	} else {
		runtime.WindowShow(a.Ctx)
	}
	if a.OnWindowHidden != nil {
		a.OnWindowHidden(hidden)
	}
}

func (a *BaseApp) SetTrayTitle(title string) {
	if goruntime.GOOS == "linux" {
		systray.SetTitle(title)
	} else if goruntime.GOOS == "windows" {
		systray.SetTooltip(title)
	}
}

func (a *BaseApp) ToggleWindowHidden() {
	a.SetWindowHidden(!a.hidden)
}
