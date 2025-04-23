package gui

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"fyne.io/systray"
	"github.com/downace/go-config"
	"github.com/downace/print-server/internal/appconfig"
	"github.com/downace/print-server/internal/common"
	"github.com/downace/print-server/internal/guiapp"
	"github.com/downace/print-server/internal/logging"
	"github.com/downace/print-server/internal/server"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"maps"
	"net"
	"net/http"
	"net/netip"
	"os"
)

//go:embed resources/window_icon.png
var windowIcon []byte

type AppTrayMenuItems struct {
	mToggleWindow *systray.MenuItem
	mToggleServer *systray.MenuItem
}

type App struct {
	appName string
	baseApp guiapp.BaseApp

	config        *config.Config[appconfig.AppConfig]
	trayMenuItems AppTrayMenuItems
	httpServer    *http.Server
}

func RunApp(appName string, assets embed.FS) error {
	app := NewApp(appName)
	return wails.Run(&options.App{
		Title:         appName,
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
			Icon: windowIcon,
		},
		Logger: logging.WailsLog,
	})
}

func NewApp(appName string) *App {
	a := App{
		appName: appName,
		config:  config.NewConfigMinimal(appconfig.NewDefaultConfig()),
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
	a.baseApp.SetTrayTitle(a.appName)

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
	if a.httpServer != nil {
		_ = a.httpServer.Shutdown(ctx)
	}
	a.baseApp.Shutdown(ctx)
}

type ServerStatus struct {
	Running     bool   `json:"running"`
	Error       string `json:"error"`
	RunningHost string `json:"runningHost"`
	RunningPort uint16 `json:"runningPort"`
}

func (a *App) GetConfig() appconfig.AppConfig {
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
			log.Printf("error getting addrs for interface %s: %s", iface.Name, err)
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

	return a.config.Transaction(func(c *appconfig.AppConfig) error {
		c.Host = newVal
		return nil
	})
}

func (a *App) UpdateServerPort(newVal uint16) error {
	if newVal == a.config.Data.Port {
		return nil
	}
	return a.config.Transaction(func(c *appconfig.AppConfig) error {
		c.Port = newVal
		return nil
	})
}

func (a *App) UpdateResponseHeaders(newVal map[string]string) error {
	if maps.Equal(newVal, a.config.Data.ResponseHeaders) {
		return nil
	}
	return a.config.Transaction(func(data *appconfig.AppConfig) error {
		data.ResponseHeaders = newVal
		return nil
	})
}

func (a *App) UpdateTLSEnabled(newVal bool) error {
	if a.config.Data.TLS.Enabled == newVal {
		return nil
	}
	return a.config.Transaction(func(data *appconfig.AppConfig) error {
		data.TLS.Enabled = newVal
		return nil
	})
}

func (a *App) UpdateTLSCertFile(newVal string) error {
	if a.config.Data.TLS.CertFile == newVal {
		return nil
	}

	if err := validateFile(newVal); err != nil {
		return err
	}

	return a.config.Transaction(func(data *appconfig.AppConfig) error {
		data.TLS.CertFile = newVal
		return nil
	})
}

func (a *App) UpdateTLSKeyFile(newVal string) error {
	if a.config.Data.TLS.KeyFile == newVal {
		return nil
	}

	if err := validateFile(newVal); err != nil {
		return err
	}

	return a.config.Transaction(func(data *appconfig.AppConfig) error {
		data.TLS.KeyFile = newVal
		return nil
	})
}

func validateFile(path string) error {
	stat, err := os.Stat(path)

	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("path is a directory")
	}
	return nil
}

func (a *App) GetServerStatus() ServerStatus {
	return a.serverStatus()
}

func (a *App) StartServer() {
	if a.httpServer != nil {
		_ = a.httpServer.Close()
	}
	a.httpServer = server.CreateServer(a.config.Data)

	go func() {
		err := server.RunServer(a.httpServer, a.config.Data)
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

func (a *App) PickFilePath() (string, error) {
	return runtime.OpenFileDialog(a.baseApp.Ctx, runtime.OpenDialogOptions{})
}

func (a *App) handleStatusChange(status ServerStatus) {
	runtime.EventsEmit(a.baseApp.Ctx, "server-status-changed", status)

	if status.Running {
		runtime.WindowSetTitle(a.baseApp.Ctx, fmt.Sprintf("%s - Running on %s:%d", a.appName, status.RunningHost, status.RunningPort))
		systray.SetIcon(trayIconRunning)
		a.baseApp.SetTrayTitle(fmt.Sprintf("%s\nRunning on %s:%d", a.appName, status.RunningHost, status.RunningPort))
		a.trayMenuItems.mToggleServer.SetTitle("Stop server")
		a.trayMenuItems.mToggleServer.Enable()
	} else {
		runtime.WindowSetTitle(a.baseApp.Ctx, fmt.Sprintf("%s - Stopped", a.appName))
		if status.Error != "" {
			systray.SetIcon(trayIconError)
		} else {
			systray.SetIcon(trayIconStopped)
		}
		a.baseApp.SetTrayTitle(fmt.Sprintf("%s\nStopped", a.appName))
		a.trayMenuItems.mToggleServer.SetTitle("Start server")
		a.trayMenuItems.mToggleServer.Enable()
	}
}
