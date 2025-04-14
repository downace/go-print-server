package logging

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
	"log"
	"os"
)

var WailsLog logger.Logger
var AppLog *log.Logger
var HttpLog *log.Logger

func InitLogs() {
	WailsLog = logger.NewFileLogger("wails.log")

	appLogFile := openLogFile("app.log")
	AppLog = log.Default()
	AppLog.SetOutput(appLogFile)

	httpLogFile := openLogFile("http.log")
	HttpLog = log.New(httpLogFile, "", log.LstdFlags)
}

func openLogFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
	if err != nil {
		panic(err)
	}

	return file
}
