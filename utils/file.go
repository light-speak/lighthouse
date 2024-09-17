package utils

import (
	"fmt"
	"os"
)

func CreateOrTruncateFile(fileName string) (*os.File, error) {
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
