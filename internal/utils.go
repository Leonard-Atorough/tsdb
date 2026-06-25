package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// getProjectRoot finds the project root by locating go.mod starting from the current working directory.
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func convertTimeToUnix(timeStr string) (int64, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}
	return parsedTime.UnixMilli(), nil
}