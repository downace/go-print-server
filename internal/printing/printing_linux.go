package printing

import (
	"github.com/samber/lo"
	"io"
	"os/exec"
	"slices"
	"strings"
)

func ListPrinters() ([]Printer, error) {
	cmd := exec.Command("lpstat", "-e")

	output, err := execAndLogCommand(cmd)

	if err != nil {
		return nil, err
	}

	return lo.Map(slices.Collect(strings.Lines(string(output))), func(line string, _ int) Printer {
		return Printer{Name: strings.TrimSpace(line)}
	}), nil
}

func PrintPDF(printer string, file io.Reader) error {
	return printPdfUsingCommand(printer, file, func(printer string, filename string) *exec.Cmd {
		return exec.Command("lp", "-d", printer, filename)
	})
}
