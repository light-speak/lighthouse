package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetModulePathFromGoMod() (string, error) {
	// 获取当前执行命令的工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	// 构造 go.mod 文件的路径
	goModPath := filepath.Join(currentDir, "go.mod")

	// 打开 go.mod 文件
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("error opening go.mod file: %v", err)
	}
	defer file.Close()

	// 读取 go.mod 文件并解析模块路径
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			// 提取模块路径
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading go.mod file: %v", err)
	}

	return "", fmt.Errorf("module path not found in go.mod")
}
