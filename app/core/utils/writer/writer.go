package writer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type CustomWriter struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *CustomWriter {
	return &CustomWriter{logger: logger}
}

func (cw CustomWriter) WriteNewFile(content, path, filename string) error {
	filePath := filepath.Join(path, filename)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)

	return err
}

func (cw CustomWriter) ReadFromFile(path string) error {
	// todo - implement
	return nil
}
