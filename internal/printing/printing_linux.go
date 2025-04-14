package printing

import (
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func ListPrinters() ([]Printer, error) {
	cmd := exec.Command("lpstat", "-e")

	log.Printf("executing %s with %q", cmd.Path, cmd.Args)

	output, err := cmd.CombinedOutput()

	log.Printf("result: %q", output)

	if err != nil {
		return nil, fmt.Errorf("%s", output)
	}

	return lo.Map(slices.Collect(strings.Lines(string(output))), func(line string, _ int) Printer {
		return Printer{Name: strings.TrimSpace(line)}
	}), nil
}

func PrintPDF(printer string, file io.Reader) error {
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

	cmd := exec.Command("lp", "-d", printer, tmpFile.Name())

	log.Printf("executing %s with %q", cmd.Path, cmd.Args)

	output, err := cmd.CombinedOutput()

	log.Printf("result: %q", output)

	if err != nil {
		return fmt.Errorf("%s", output)
	}
	return nil
}
