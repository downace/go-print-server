package printing

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type Printer struct {
	Name string `json:"name"`
}

var ErrNotSupported = fmt.Errorf("method not supported on %s", runtime.GOOS)
var ErrRequestError = fmt.Errorf("request error")

var browserPage *rod.Page

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

func PrintFromUrl(printer string, url string, options *proto.PagePrintToPDF) error {
	pdfFile, err := urlToPdf(url, options)

	if err != nil {
		return err
	}

	return PrintPDF(printer, pdfFile)
}

func initBrowserPage() (*rod.Page, error) {
	if browserPage != nil {
		return browserPage, nil
	}
	browser := rod.New()

	err := browser.Connect()

	if err != nil {
		return nil, err
	}

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}

	browserPage = page
	return page, nil
}

func urlToPdf(url string, options *proto.PagePrintToPDF) (pdfFile io.Reader, err error) {
	page, err := initBrowserPage()
	if err != nil {
		return
	}

	responseReceivedEvent := proto.NetworkResponseReceived{}
	waitResponse := page.WaitEvent(&responseReceivedEvent)
	err = page.Navigate(url)
	if err != nil {
		return
	}
	waitResponse()

	if resp := responseReceivedEvent.Response; resp.Status >= 300 {
		err = fmt.Errorf("%w: response from URL was %d %s", ErrRequestError, resp.Status, resp.StatusText)
		return
	}

	err = page.WaitLoad()
	if err != nil {
		return
	}

	pdfFile, err = page.PDF(options)
	return pdfFile, err
}

func execAndLogCommand(cmd *exec.Cmd) (output []byte, err error) {
	log.Printf("executing %s with %q", cmd.Path, cmd.Args)

	output, err = cmd.CombinedOutput()

	log.Printf("result: %q", output)

	if err != nil {
		return nil, fmt.Errorf("%s", output)
	}

	return output, nil
}

func printPdfUsingCommand(printer string, file io.Reader, commandFactory func(printer string, filename string) *exec.Cmd) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "print-server-*.pdf")

	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)

	if err != nil {
		return err
	}

	cmd := commandFactory(printer, tmpFile.Name())

	_, err = execAndLogCommand(cmd)

	return err
}
