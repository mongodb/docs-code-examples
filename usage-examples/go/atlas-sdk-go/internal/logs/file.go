package logs

import (
	"fmt"
	"io"
	"os"

	"atlas-sdk-go/internal"
)

// WriteToFile copies everything from r into a new file at path.
// It will create or truncate that file.
func WriteToFile(r io.Reader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer internal.SafeClose(f)

	if err := internal.SafeCopy(f, r); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
