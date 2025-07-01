package fileutils

import (
	"fmt"
	"path/filepath"
	"time"
)

// GenerateOutputPath creates a file path based on the base directory, prefix, and optional extension.
func GenerateOutputPath(baseDir, prefix, extension string) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("baseDir cannot be empty")
	}
	if prefix == "" {
		return "", fmt.Errorf("prefix cannot be empty")
	}

	timestamp := time.Now().Format("20060102")
	var filename string
	if extension == "" {
		filename = fmt.Sprintf("%s_%s", prefix, timestamp)
	} else {
		filename = fmt.Sprintf("%s_%s.%s", prefix, timestamp, extension)
	}

	return filepath.Join(baseDir, filename), nil
}
