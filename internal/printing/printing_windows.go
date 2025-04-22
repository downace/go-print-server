package printing

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/downace/print-server/internal/common"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
)

//go:embed SumatraPDF.exe
var embedFs embed.FS

func ListPrinters() ([]Printer, error) {
	cmd := exec.Command("wmic", "printer", "list", "brief", "/format:csv")

	log.Printf("executing %s with %q", cmd.Path, cmd.Args)

	output, err := cmd.CombinedOutput()

	log.Printf("result: %q", output)

	if err != nil {
		return nil, fmt.Errorf("%s", output)
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

	cmd := exec.Command(sumatra, "-print-to", printer, "-silent", tmpFile.Name())

	log.Printf("executing %s with %q", cmd.Path, cmd.Args)

	output, err := cmd.CombinedOutput()

	log.Printf("result: %q", output)

	if err != nil {
		return fmt.Errorf("%s", output)
	}
	return nil
}
