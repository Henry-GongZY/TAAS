package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const (
	// PATH mode: compute MD5 from a file path.
	PATH = iota
	// BUFFER mode: compute MD5 from a byte slice in memory.
	BUFFER
)

// getFileMD5 computes the MD5 hash of data based on the specified mode.
// The data parameter can be a file path (string) or a byte slice ([]byte).
func getFileMD5(data interface{}, mode int) (string, error) {
	var reader io.Reader

	switch mode {
	case PATH:
		// Ensure data is a string file path.
		filePath, ok := data.(string)
		if !ok {
			return "", fmt.Errorf("in PATH mode, data must be a string (file path)")
		}

		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file) // Ensure the file is closed.

		reader = file

	case BUFFER:
		// Ensure data is a byte slice.
		buffer, ok := data.([]byte)
		if !ok {
			return "", fmt.Errorf("in BUFFER mode, data must be a byte slice")
		}
		reader = bytes.NewReader(buffer)

	default:
		return "", fmt.Errorf("invalid mode: %d", mode)
	}

	// Calculate MD5 hash using a single logic for all modes.
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", fmt.Errorf("failed to copy data for MD5 calculation: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
