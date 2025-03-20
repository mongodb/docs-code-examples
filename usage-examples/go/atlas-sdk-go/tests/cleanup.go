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
	projectRoot := "./" // Project root directory
	files, err := filepath.Glob(filepath.Join(projectRoot, "*.gz"))
	if err != nil {
		return fmt.Errorf("failed to find .gz files: %w", err)
	}
	for _, file := range files {
		if err := deleteFile(file); err != nil {
			log.Printf("Failed to delete file %s: %v", file, err)
		}
	}
	return nil
}
