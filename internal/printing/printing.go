package printing

import (
	"fmt"
	"runtime"
)

type Printer struct {
	Name string `json:"name"`
}

var ErrNotSupported = fmt.Errorf("method not supported on %s", runtime.GOOS)
