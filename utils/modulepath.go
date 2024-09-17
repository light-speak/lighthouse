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
		return "", fmt.Errorf("获取当前目录失败: %v", err)
	}

	// 构造 go.mod 文件的路径
	goModPath := filepath.Join(currentDir, "go.mod")

	// 打开并读取 go.mod 文件
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("打开 go.mod 文件失败: %v", err)
	}
	defer file.Close()

	// 使用 bufio.Scanner 高效读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			// 提取并返回模块路径
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取 go.mod 文件失败: %v", err)
	}

	return "", fmt.Errorf("在 go.mod 文件中未找到模块路径")
}
