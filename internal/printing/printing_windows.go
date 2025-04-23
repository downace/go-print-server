package printing

import (
	"bytes"
	"embed"
	"encoding/csv"
	"github.com/downace/print-server/internal/common"
	"io"
	"os/exec"
	"slices"
)

//go:embed SumatraPDF.exe
var embedFs embed.FS

func ListPrinters() ([]Printer, error) {
	cmd := exec.Command("wmic", "printer", "list", "brief", "/format:csv")

	output, err := execAndLogCommand(cmd)

	if err != nil {
		return nil, err
	}

	c := csv.NewReader(common.NewNormalizedLinesReader(bytes.NewReader(output)))
	c.ReuseRecord = true
	c.TrimLeadingSpace = true

	var headers []string = nil
	var printers []Printer

	for {
		record, err := c.Read()
		if err != nil {
			break
		}
		if headers == nil {
			headers = make([]string, len(record))
			copy(headers, record)
		} else {
			namePos := slices.Index(headers, "Name")
			name := record[namePos]
			printers = append(printers, Printer{Name: name})
		}
	}

	return printers, nil
}

func PrintPDF(printer string, file io.Reader) error {
	sumatra, err := common.MaterializeEmbeddedFile(embedFs, "SumatraPDF.exe")

	if err != nil {
		return err
	}

	return printPdfUsingCommand(printer, file, func(printer string, filename string) *exec.Cmd {
		return exec.Command(sumatra, "-print-to", printer, "-silent", filename)
	})
}
