package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/downace/go-config"
	"github.com/downace/print-server/internal/appconfig"
	"github.com/downace/print-server/internal/server"
	"github.com/samber/lo"
	"github.com/ttacon/chalk"
	"gopkg.in/yaml.v3"
	"math"
	"net/http"
	"net/netip"
	"strings"
)

type headersFlag map[string]string

func (f *headersFlag) String() string {
	return fmt.Sprint(*f)
}

func (f *headersFlag) Set(value string) error {
	pair := strings.SplitN(value, ":", 2)
	if len(pair) < 2 {
		return fmt.Errorf("invalid header: %q", value)
	}
	name := strings.TrimSpace(pair[0])
	val := strings.TrimSpace(pair[1])
	if val == "" {
		return fmt.Errorf("invalid header: %q", value)
	}
	if *f == nil {
		*f = make(headersFlag)
	}
	(*f)[name] = val
	return nil
}

func RunApp(args []string) error {
	conf := config.NewConfigMinimal(appconfig.NewDefaultConfig())
	lo.Must0(conf.Load())

	err := setConfigFromArgs(&conf.Data, args)

	if err != nil {
		return err
	}

	fmt.Printf("Using config:\n\n%s\n\n", chalk.Yellow.Color(string(lo.Must(yaml.Marshal(conf.Data)))))
	fmt.Printf("Edit %sconfig.yaml%s file or use CLI flags (see %scli -help%s) to adjust config%s\n",
		chalk.Magenta,
		chalk.Blue,
		chalk.Magenta,
		chalk.Blue,
		chalk.Reset,
	)

	serv := server.CreateServer(conf.Data)

	var proto string
	if conf.Data.TLS.Enabled {
		proto = "https"
	} else {
		proto = "http"
	}

	fmt.Println()
	fmt.Println(chalk.Green.Color(fmt.Sprintf("Running server on %s://%s:%d", proto, conf.Data.Host, conf.Data.Port)))

	err = server.RunServer(serv, conf.Data)

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func setConfigFromArgs(conf *appconfig.AppConfig, args []string) error {
	fs := flag.NewFlagSet("cli", flag.ExitOnError)

	var host string
	var port int
	var respHeaders headersFlag
	var enableTls bool
	var certFile string
	var keyFile string

	fs.StringVar(&host, "host", "", "listen host")
	fs.IntVar(&port, "port", 0, "listen port")
	fs.Var(&respHeaders, "header", "response header, can be specified multiple times")
	fs.BoolVar(&enableTls, "tls", false, "enable TLS")
	fs.StringVar(&certFile, "cert-file", "", "TLS certificate file path")
	fs.StringVar(&keyFile, "key-file", "", "TLS key file path")

	err := fs.Parse(args)

	if err != nil {
		return err
	}

	var validationError error

	fs.Visit(func(f *flag.Flag) {
		if validationError != nil {
			return
		}
		switch f.Name {
		case "host":
			_, e := netip.ParseAddr(host)
			if e != nil {
				validationError = fmt.Errorf("invalid host: %w", e)
			} else {
				conf.Host = host
			}
		case "port":
			if port < 0 || port > math.MaxUint16 {
				validationError = fmt.Errorf("invalid port: %d", port)
			} else {
				conf.Port = uint16(port)
			}
		case "header":
			conf.ResponseHeaders = respHeaders
		case "tls":
			conf.TLS.Enabled = enableTls
		case "cert-file":
			conf.TLS.CertFile = certFile
		case "key-file":
			conf.TLS.KeyFile = keyFile
		}
	})

	return validationError
}
