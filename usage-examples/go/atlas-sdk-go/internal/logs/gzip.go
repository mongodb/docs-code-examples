package logs

import (
	"compress/gzip"
	"fmt"
	"os"

	"atlas-sdk-go/internal"
)

// DecompressGzip opens a .gz file and unpacks to specified destination.
func DecompressGzip(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", srcPath, err)
	}
	defer internal.SafeClose(srcFile)

	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("gzip reader %s: %w", srcPath, err)
	}
	defer internal.SafeClose(gzReader)

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", destPath, err)
	}
	defer internal.SafeClose(destFile)

	if err := internal.SafeCopy(destFile, gzReader); err != nil {
		return fmt.Errorf("decompress to %s: %w", destPath, err)
	}
	return nil
}
