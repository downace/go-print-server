package common

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ListenChannel[T any](channel chan T, cb func(value T)) {
	go func() {
		for value := range channel {
			cb(value)
		}
	}()
}

func copyEmbeddedFile(realPath string, embedPath string, embedFs embed.FS) error {
	embeddedFile, err := embedFs.Open(embedPath)
	if err != nil {
		return err
	}
	defer embeddedFile.Close()
	realFile, err := os.Create(realPath)
	if err != nil {
		return err
	}
	defer realFile.Close()
	_, err = realFile.ReadFrom(embeddedFile)

	return err
}

func MaterializeEmbeddedFile(embedFs embed.FS, path string) (string, error) {
	execPath, err := os.Executable()

	if err != nil {
		return "", err
	}

	realPath := filepath.Join(filepath.Dir(execPath), path)
	_, err = os.Stat(realPath)

	if errors.Is(err, os.ErrNotExist) {
		err = copyEmbeddedFile(realPath, path, embedFs)
	}
	if err != nil {
		return "", err
	}
	return realPath, nil
}

type NormalizedLinesReader struct {
	scanner *bufio.Scanner
	buffer  *bytes.Buffer
}

func NewNormalizedLinesReader(reader io.Reader) *NormalizedLinesReader {
	return &NormalizedLinesReader{
		scanner: bufio.NewScanner(reader),
		buffer:  &bytes.Buffer{},
	}
}

func (r *NormalizedLinesReader) Read(p []byte) (int, error) {
	n, err := r.buffer.Read(p)

	if err != nil && err != io.EOF {
		return n, err
	}

	if n == len(p) {
		return n, nil
	}

	if r.scanner.Scan() {
		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			return n, r.scanner.Err()
		}

		lineBytes := []byte(line + "\n")
		for _, b := range lineBytes {
			if n == len(p) {
				r.buffer.WriteByte(b)
			} else {
				p[n] = b
				n++
			}
		}
		err = nil
	} else {
		err = r.scanner.Err()
		if err == nil {
			err = io.EOF
		}
	}
	return n, err
}
