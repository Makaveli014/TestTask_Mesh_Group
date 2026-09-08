package model

import (
	"os"
	"path/filepath"
	"time"
)

// CleanupOldLogs removes files in the specified directory that are older than maxAgeDays.
func CleanupOldLogs(logDir string, maxAgeDays int) ([]string, error) {
	var deletedFiles []string
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cutoff := now.AddDate(0, 0, -maxAgeDays)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			err := os.Remove(filepath.Join(logDir, f.Name()))
			if err == nil {
				deletedFiles = append(deletedFiles, f.Name())
			}
		}
	}
	return deletedFiles, nil
}
