package printing

import (
	"fmt"
)

func ListPrinters() ([]Printer, error) {
	return nil, fmt.Errorf("ListPrinters: %w", ErrNotSupported)
}

func PrintPDF(file os.File) error {
	return nil, fmt.Errorf("PrintPDF: %w", ErrNotSupported)
}
