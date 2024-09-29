package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectPath() (string, error) {
	modPath, err := GetModPath(nil)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(modPath)
	return dir, nil
}

// GetModPath Get the path to the go.mod file
func GetModPath(path *string) (string, error) {
	var currentDir string
	var err error
	if path == nil {
		currentDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	} else {
		currentDir = *path
	}

	for {
		modPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return modPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("go.mod file not found")
}

// GetCallerPath Get the path of the caller
func GetCallerPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}
	return filepath.Dir(currentFilePath), nil
}
