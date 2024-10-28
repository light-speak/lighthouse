package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetProjectPath get the project path
func GetProjectPath() (string, error) {
	modPath, err := GetModPath(nil)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(modPath)
	return dir, nil
}

// GetModPath get the path location of the go.mod file
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

// GetPackageFromPath Get the package name from the path ending
func GetPackageFromPath(path string) (string, error) {
	paths := strings.Split(path, ",")
	if len(paths) == 0 {
		return "", fmt.Errorf("path is not valid")
	}
	return paths[len(paths)-1], nil
}

// EnsureDir makes sure the specified directory exists
// If the directory doesn't exist, it creates it with full permissions
func EnsureDir(path string) error {
	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}


// getModuleName reads the go.mod file and returns the module name
func GetModuleName(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			// Extract the module name
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("module name not found")
}