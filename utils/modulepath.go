package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetModulePath() (string, error) {
	// 获取当前执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取执行文件路径失败: %v", err)
	}

	// 获取执行文件所在目录
	execDir := filepath.Dir(execPath)

	// 尝试从go.mod文件获取模块路径
	modulePath, err := getModulePathFromGoMod(execDir)
	if err == nil {
		return modulePath, nil
	}

	// 如果无法从go.mod获取，尝试从可执行文件名推断
	return inferModulePathFromExecutable(execPath), nil
}

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

func inferModulePathFromExecutable(execPath string) string {
	// 获取可执行文件名（不包含扩展名）
	baseName := filepath.Base(execPath)
	if runtime.GOOS == "windows" {
		baseName = strings.TrimSuffix(baseName, ".exe")
	}

	// 假设模块路径与可执行文件名相同
	return baseName
}
