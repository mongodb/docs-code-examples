package fileutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// WriteToFile copies everything from r into a new file at path.
// It will create or truncate that file.
func WriteToFile(r io.Reader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer SafeClose(f)

	if err := SafeCopy(f, r); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// SafeClose closes c and logs a warning on error.
func SafeClose(c io.Closer) {
	if c != nil {
		if err := c.Close(); err != nil {
			log.Printf("warning: close failed: %v", err)
		}
	}
}

// SafeCopy copies src → dst and propagates any error (after logging).
func SafeCopy(dst io.Writer, src io.Reader) error {
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("warning: copy failed: %v", err)
		return err
	}
	return nil
}

// :remove-start:

// SafeDelete removes files generated in the specified directory.
// NOTE: INTERNAL ONLY FUNCTION
func SafeDelete(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if removeErr := os.Remove(path); removeErr != nil {
				log.Printf("warning: failed to delete file %s: %v", path, removeErr)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// :remove-end:
