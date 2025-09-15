package typeutils

import (
	"fmt"
	"strings"
)

// FirstNonEmpty returns the first non-empty, non-whitespace string from values, or an empty string if none found.
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// ParseBool interprets a string as a boolean. It returns true for "true", "1", "yes", "y" (case insensitive, trimmed), and false otherwise.
func ParseBool(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "true" || v == "1" || v == "yes" || v == "y"
}

// SuffixZeroed returns a formatted string if zeroed > 0, otherwise an empty string.
func SuffixZeroed(zeroed int, region string) string {
	if zeroed == 0 {
		return ""
	}
	return fmt.Sprintf(", zeroed electable nodes in region %s", region)
}

// DefaultIfBlank returns d if v is an empty string, otherwise returns v.
func DefaultIfBlank(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
