package model

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanupOldLogs(t *testing.T) {
	// Setup temporary directory for logs
	tmpDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new file
	newFile := filepath.Join(tmpDir, "new.log")
	if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create an old file
	oldFile := filepath.Join(tmpDir, "old.log")
	if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	// Manually set the modification time to 10 days ago
	oldTime := time.Now().AddDate(0, 0, -10)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	// Run cleanup with maxAgeDays = 7
	deleted, err := CleanupOldLogs(tmpDir, 7)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that old.log was deleted
	foundOld := false
	for _, f := range deleted {
		if f == "old.log" {
			foundOld = true
			break
		}
	}
	if !foundOld {
		t.Error("Expected old.log to be deleted, but it wasn't")
	}

	// Check that new.log still exists
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("Expected new.log to persist, but it was deleted")
	}
}
