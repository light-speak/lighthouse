package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateOrTruncateFile(fileName string) (*os.File, error) {
	// 获取文件所在的目录
	dir := filepath.Dir(fileName)

	// 如果目录不存在，则创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录时出错: %v", err)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("打开或创建文件时出错: %v", err)
	}

	// 设置文件偏移量为起始位置
	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("设置文件偏移量时出错: %v", err)
	}

	return file, nil
}


