// :snippet-start: utility-functions-full-example
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LoadEnvWithDefault loads an environment variable or returns a default value if not set.
func LoadEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// HandleError logs the error and exits the program.
func HandleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

// FormatResponseAsJSON reads an io.ReadCloser, formats it as pretty JSON, and returns it as a string.
func FormatResponseAsJSON(body io.ReadCloser) (string, error) {
	if body == nil {
		return "", fmt.Errorf("response body is nil")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			HandleError(err, "Error closing body")
		}
	}(body)

	// Read the body into a buffer
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert raw response to JSON
	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, buf.Bytes(), "", "  ") // Pretty-print JSON
	if err != nil {
		return "", fmt.Errorf("failed to format response as JSON: %w", err)
	}

	return formattedJSON.String(), nil
}

// IsEmptyString checks if a string is empty or contains only whitespace.
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Retry retries a function up to a specified number of times with a delay between attempts.
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			if i >= (attempts - 1) {
				return err
			}
			time.Sleep(sleep)
			continue
		}
		return nil
	}
	return fmt.Errorf("reached maximum retry attempts")
}

// :snippet-end: [utility-functions-full-example]
