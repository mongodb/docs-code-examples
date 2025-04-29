package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// deleteFile deletes a specified file during cleanup after tests.
func deleteFile(fileName string) error {
	err := os.Remove(fileName)
	log.Printf("Deleting file %s", fileName)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fileName, err)
	}
	return nil
}

// CleanupFiles deletes generated *.gz and *.log files from the project root directory
func CleanupFiles() error {
	projectRoot := "./"
	var filesToDelete []string

	for _, pattern := range []string{"*.gz", "*.log"} {
		files, err := filepath.Glob(filepath.Join(projectRoot, pattern))
		if err != nil {
			return fmt.Errorf("failed to glob %q: %w", pattern, err)
		}
		filesToDelete = append(filesToDelete, files...)
	}

	for _, file := range filesToDelete {
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to delete file %s: %v", file, err)
		} else {
			log.Printf("Deleted %s", file)
		}
	}

	return nil
}
