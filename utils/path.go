package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

// GetGoModPath returns the path to go.mod file using go env GOMOD
func GetGoModPath() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run 'go env GOMOD': %w", err)
	}

	gomodPath := strings.TrimSpace(string(output))
	if gomodPath == "" || gomodPath == "/dev/null" {
		return "", errors.New("not in a Go module (go env GOMOD returned empty or /dev/null)")
	}

	return gomodPath, nil
}

// GetModuleRoot returns the module root directory (directory containing go.mod)
func GetModuleRoot() (string, error) {
	gomodPath, err := GetGoModPath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(gomodPath), nil
}

// GetModPath returns the module path using go env GOMOD
// Note: projectPath parameter is deprecated and ignored; go env GOMOD is used instead
func GetModPath(projectPath *string) (string, error) {
	gomodPath, err := GetGoModPath()
	if err != nil {
		return "", err
	}

	// Read go.mod file from the path returned by go env GOMOD
	data, err := os.ReadFile(gomodPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod at %s: %w", gomodPath, err)
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
// Note: projectPath parameter is deprecated and ignored; go env GOMOD is used instead
func GetPkgPath(projectPath, filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("file path is empty")
	}

	// Get module path name
	modPath, err := GetModPath(nil)
	if err != nil {
		return "", err
	}

	// Get module root directory
	modRoot, err := GetModuleRoot()
	if err != nil {
		return "", err
	}

	// Convert filePath to absolute path if it's relative
	absFilePath := filePath
	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		absFilePath = filepath.Join(cwd, filePath)
	}

	// Get directory of file
	dir := filepath.Dir(absFilePath)

	// Calculate relative path from module root
	relPath, err := filepath.Rel(modRoot, dir)
	if err != nil {
		return "", fmt.Errorf("failed to calculate relative path: %w", err)
	}

	if relPath == "." {
		return modPath, nil
	}

	// Combine full package path
	return filepath.Join(modPath, relPath), nil
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
