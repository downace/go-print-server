package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"fyne.io/systray"
	"github.com/downace/go-config"
	"github.com/downace/print-server/internal/app"
	"github.com/downace/print-server/internal/common"
	"github.com/downace/print-server/internal/server"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net"
	"net/http"
	"net/netip"
)

//go:embed build/trayicon_default.png
var trayIconDefault []byte

//go:embed build/trayicon_running.png
var trayIconRunning []byte

//go:embed build/trayicon_stopped.png
var trayIconStopped []byte

//go:embed build/trayicon_error.png
var trayIconError []byte

type AppTrayMenuItems struct {
	mToggleWindow *systray.MenuItem
	mToggleServer *systray.MenuItem
}

type AppConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint16 `yaml:"port" json:"port"`
}

// App struct
type App struct {
	baseApp app.BaseApp

	config        *config.Config[AppConfig]
	trayMenuItems AppTrayMenuItems
	httpServer    *http.Server
}

func NewApp() *App {
	a := App{
		config: config.NewConfigMinimal(AppConfig{
			Host: "0.0.0.0",
			Port: 8888,
		}),
	}
	a.baseApp.InitTray = a.initTray
	a.baseApp.OnWindowHidden = a.onWindowHidden
	return &a
}

func (a *App) startup(ctx context.Context) {
	a.baseApp.Startup(ctx)

	lo.Must0(a.config.Load())
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return a.baseApp.BeforeClose(ctx)
}

func (a *App) onWindowHidden(hidden bool) {
	if hidden {
		a.trayMenuItems.mToggleWindow.SetTitle("Show window")
	} else {
		a.trayMenuItems.mToggleWindow.SetTitle("Hide window")
	}
}

func (a *App) initTray() {
	systray.SetIcon(trayIconDefault)
	systray.SetTitle(AppName)

	a.trayMenuItems.mToggleWindow = systray.AddMenuItem("Hide window", "Hide window")
	systray.AddSeparator()
	a.trayMenuItems.mToggleServer = systray.AddMenuItem("Start server", "Start server")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit")

	common.ListenChannel(a.trayMenuItems.mToggleWindow.ClickedCh, func(_ struct{}) {
		a.baseApp.ToggleWindowHidden()
	})

	common.ListenChannel(a.trayMenuItems.mToggleServer.ClickedCh, func(_ struct{}) {
		if a.httpServer != nil {
			a.StopServer()
		} else {
			a.StartServer()
		}
	})

	common.ListenChannel(mQuit.ClickedCh, func(_ struct{}) {
		a.baseApp.Quit()
	})

	a.handleStatusChange(ServerStatus{Running: false})
}

func (a *App) shutdown(ctx context.Context) {
	_ = a.httpServer.Shutdown(ctx)
	a.baseApp.Shutdown(ctx)
}

type ServerStatus struct {
	Running     bool   `json:"running"`
	Error       string `json:"error"`
	RunningHost string `json:"runningHost"`
	RunningPort uint16 `json:"runningPort"`
}

func (a *App) GetConfig() AppConfig {
	return a.config.Data
}

type NetInterface struct {
	Name string `json:"name"`
	IsUp bool   `json:"isUp"`
}

type NetInterfaceAddress struct {
	Ip        string       `json:"ip"`
	Interface NetInterface `json:"interface"`
}

func (a *App) GetAvailableAddrs() ([]NetInterfaceAddress, error) {
	ips := make([]NetInterfaceAddress, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()

		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch ipAddr := addr.(type) {
			case *net.IPNet:
				ip = ipAddr.IP
			case *net.IPAddr:
				ip = ipAddr.IP
			}
			ips = append(ips, NetInterfaceAddress{
				Ip: ip.String(),
				Interface: NetInterface{
					Name: iface.Name,
					IsUp: iface.Flags&net.FlagUp != 0,
				},
			})
		}
	}

	return ips, nil
}

func (a *App) serverStatus() ServerStatus {
	if a.httpServer == nil {
		return ServerStatus{Running: false}
	}

	addrPort := netip.MustParseAddrPort(a.httpServer.Addr)

	return ServerStatus{
		Running:     true,
		RunningHost: addrPort.Addr().String(),
		RunningPort: addrPort.Port(),
	}
}

func (a *App) UpdateServerHost(newVal string) error {
	if newVal == a.config.Data.Host {
		return nil
	}
	_, err := netip.ParseAddr(newVal)
	if err != nil {
		return err
	}

	return a.config.Transaction(func(c *AppConfig) error {
		c.Host = newVal
		return nil
	})
}

func (a *App) UpdateServerPort(newVal uint16) error {
	if newVal == a.config.Data.Port {
		return nil
	}
	return a.config.Transaction(func(c *AppConfig) error {
		c.Port = newVal
		return nil
	})
}

func (a *App) GetServerStatus() ServerStatus {
	return a.serverStatus()
}

func (a *App) StartServer() {
	if a.httpServer != nil {
		_ = a.httpServer.Close()
	}
	a.httpServer = server.CreateServer(netip.AddrPortFrom(netip.MustParseAddr(a.config.Data.Host), a.config.Data.Port))

	go func() {
		err := a.httpServer.ListenAndServe()
		a.httpServer = nil
		if errors.Is(err, http.ErrServerClosed) {
			a.handleStatusChange(ServerStatus{Running: false})
		} else {
			a.handleStatusChange(ServerStatus{Running: false, Error: err.Error()})
		}
	}()

	a.handleStatusChange(a.serverStatus())
}

func (a *App) StopServer() {
	if a.httpServer != nil {
		_ = a.httpServer.Close()
	}
}

func (a *App) handleStatusChange(status ServerStatus) {
	runtime.EventsEmit(a.baseApp.Ctx, "server-status-changed", status)

	if status.Running {
		runtime.WindowSetTitle(a.baseApp.Ctx, fmt.Sprintf("%s - Running on %s:%d", AppName, status.RunningHost, status.RunningPort))
		systray.SetIcon(trayIconRunning)
		systray.SetTitle(fmt.Sprintf("%s\nRunning on %s:%d", AppName, status.RunningHost, status.RunningPort))
		a.trayMenuItems.mToggleServer.SetTitle("Stop server")
		a.trayMenuItems.mToggleServer.Enable()
	} else {
		runtime.WindowSetTitle(a.baseApp.Ctx, fmt.Sprintf("%s - Stopped", AppName))
		if status.Error != "" {
			systray.SetIcon(trayIconError)
		} else {
			systray.SetIcon(trayIconStopped)
		}
		systray.SetTitle(fmt.Sprintf("%s\nStopped", AppName))
		a.trayMenuItems.mToggleServer.SetTitle("Start server")
		a.trayMenuItems.mToggleServer.Enable()
	}
}
