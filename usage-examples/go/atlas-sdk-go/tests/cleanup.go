package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// deleteFile deletes a specified file
func deleteFile(fileName string) error {
	err := os.Remove(fileName)
	log.Printf("Deleting file %s", fileName)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fileName, err)
	}
	return nil
}

// CleanupGzFiles deletes all .gz files in the project root directory
func CleanupGzFiles() error {
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
