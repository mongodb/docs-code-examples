package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Utility functions
func FirstNonEmpty(cli, config string) string {
	if cli != "" {
		return cli
	}
	return config
}
func FirstNonEmptyArray(cli, config []string) []string {
	if len(cli) > 0 {
		return cli
	}
	return config
}

func FirstNonZero(cli, config int) int {
	if cli != 0 {
		return cli
	}
	return config
}

// FormatResponseAsJSON reads an io.ReadCloser, formats it as pretty JSON, and returns it as a string.
func FormatResponseAsJSON(body io.ReadCloser) (string, error) {
	if body == nil {
		return "", fmt.Errorf("response body is nil")
	}
	defer body.Close()

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
