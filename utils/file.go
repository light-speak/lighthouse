package utils

import (
	"fmt"
	"os"
)

func CreateOrTruncateFile(fileName string) (*os.File, error) {
	var file *os.File
	var err error

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %v", err)
		}
	} else {
		file, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
	}

	// 清空文件内容
	err = file.Truncate(0)
	if err != nil {
		return nil, fmt.Errorf("error truncating file: %v", err)
	}

	// 设置文件偏移量为起始位置
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking file: %v", err)
	}

	return file, nil
}
