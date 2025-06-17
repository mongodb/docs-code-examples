package internal

import (
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"log"
	"os"            // :remove:
	"path/filepath" // :remove:
)

// FormatAPIError formats an error returned by the Atlas API with additional context.
func FormatAPIError(operation string, params interface{}, err error) error {
	if apiErr, ok := admin.AsError(err); ok && apiErr.GetDetail() != "" {
		return fmt.Errorf("%s %v: %w: %s", operation, params, err, apiErr.GetDetail())
	}
	return fmt.Errorf("%s %v: %w", operation, params, err)
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
