package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetModulePath 获取模块路径
func GetModulePath() (string, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 尝试从当前目录的go.mod文件获取模块路径
	modulePath, err := getModulePathFromGoMod(currentDir)
	if err == nil {
		return modulePath, nil
	}

	// 如果在当前目录找不到go.mod，尝试向上查找父目录
	for dir := currentDir; dir != "/"; dir = filepath.Dir(dir) {
		modulePath, err := getModulePathFromGoMod(dir)
		if err == nil {
			return modulePath, nil
		}
	}

	// 如果无法找到go.mod，返回错误
	return "", fmt.Errorf("无法找到项目的go.mod文件")
}

// getModulePathFromGoMod 从go.mod文件中获取模块路径
func getModulePathFromGoMod(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("读取go.mod文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("在go.mod文件中未找到模块路径")
}


func GetLibraryPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}
	return filepath.Dir(filepath.Dir(currentFilePath)), nil
}
