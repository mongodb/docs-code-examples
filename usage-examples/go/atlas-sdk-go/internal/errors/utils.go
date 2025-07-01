package errors

import (
	"fmt"
	"log"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FormatError formats an error message for a specific operation and entity ID.
func FormatError(operation string, entityID string, err error) error {
	if apiErr, ok := admin.AsError(err); ok && apiErr.GetDetail() != "" {
		return fmt.Errorf("%s for %s: %w: %s", operation, entityID, err, apiErr.GetDetail())
	}
	return fmt.Errorf("%s for %s: %w", operation, entityID, err)
}

// WithContext adds context information to an error
func WithContext(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

// ExitWithError prints an error message with context and exits the program.
func ExitWithError(context string, err error) {
	log.Fatalf("%s: %v", context, err)
	// Note: log.Fatalf calls os.Exit(1)
}
