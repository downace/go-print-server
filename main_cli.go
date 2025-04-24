//go:build cli

package main

import (
	"github.com/downace/print-server/internal/cli"
	"github.com/downace/print-server/internal/logging"
	"os"
)

func main() {
	logging.InitLogs()

	err := cli.RunApp()

	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}
