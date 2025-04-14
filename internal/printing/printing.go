package printing

import (
	"fmt"
	"net/http"
	"runtime"
)

type Printer struct {
	Name string `json:"name"`
}

var ErrNotSupported = fmt.Errorf("method not supported on %s", runtime.GOOS)
var ErrRequestError = fmt.Errorf("request error")

func PrintPDFFromUrl(printer string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: response from URL was %s", ErrRequestError, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return fmt.Errorf("%w: downloaded file is %s, expected %s", ErrRequestError, contentType, "application/pdf")
	}

	return PrintPDF(printer, resp.Body)
}
