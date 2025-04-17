package cli

import (
	"errors"
	"fmt"
	"github.com/downace/go-config"
	"github.com/downace/print-server/internal/appconfig"
	"github.com/downace/print-server/internal/server"
	"github.com/samber/lo"
	"github.com/ttacon/chalk"
	"gopkg.in/yaml.v3"
	"net/http"
)

func RunApp() error {
	conf := config.NewConfigMinimal(appconfig.NewDefaultConfig())
	lo.Must0(conf.Load())

	fmt.Printf("Using config:\n\n%s\n\n", chalk.Yellow.Color(string(lo.Must(yaml.Marshal(conf.Data)))))
	fmt.Println(chalk.Magenta.Color("You can adjust config by creating or modifying config.yaml file"))

	serv := server.CreateServer(conf.Data)

	var proto string
	if conf.Data.TLS.Enabled {
		proto = "https"
	} else {
		proto = "http"
	}

	fmt.Println()
	fmt.Println(chalk.Green.Color(fmt.Sprintf("Running server on %s://%s:%d", proto, conf.Data.Host, conf.Data.Port)))

	err := server.RunServer(serv, conf.Data)

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
