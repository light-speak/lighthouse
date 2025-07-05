package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// GetProjectPath returns the project root path
func GetProjectPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetModPath returns the module path
func GetModPath(projectPath *string) (string, error) {
	path := "."
	if projectPath != nil {
		path = *projectPath
	}

	// Read go.mod file
	data, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return "", err
	}

	// Parse module path
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", errors.New("module path not found in go.mod")
}

// GetPkgPath returns the package path
func GetPkgPath(projectPath, filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("file path is empty")
	}

	// Get directory of file
	dir := filepath.Dir(filePath)
	if dir == "." {
		return "", nil
	}

	// Get module path
	modPath, err := GetModPath(&projectPath)
	if err != nil {
		return "", err
	}

	// Combine full package path
	return filepath.Join(modPath, dir), nil
}

// GetGoPath returns GOPATH
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// If GOPATH not set in env, return default
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "go")
	}
	return gopath
}

// GetFilePath returns the file path
func GetFilePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}

	// 先检查路径是否存在
	fileInfo, err := os.Stat(path)
	if err == nil {
		// 如果路径存在，检查是否是目录
		if fileInfo.IsDir() {
			return "", errors.New("path is a directory")
		}
		return filepath.Base(path), nil
	}

	// 如果路径不存在，尝试判断是否看起来像一个目录
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return "", errors.New("path is a directory")
	}

	// 检查最后一个部分是否包含扩展名
	base := filepath.Base(path)
	if !strings.Contains(base, ".") {
		return "", errors.New("path is a directory")
	}

	return base, nil
}

// GetFileDir returns the directory containing the file
func GetFileDir(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}

	dir := filepath.Dir(path)
	if dir == "." {
		return ".", nil
	}

	return dir, nil
}

// MkdirAll creates a directory named path, along with any necessary parents
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
