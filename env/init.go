package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

var (
	loadedEnv = false
)

func init() {
	// 获取当前工作目录（即程序执行的目录）
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current working directory: %v\n", err)
		return
	}

	// 构造.env文件的路径
	envFilePath := filepath.Join(currentDir, ".env")
	if err := initEnv(&envFilePath); err != nil {
		fmt.Printf("Failed to initialize environment: %v\n", err)
	}
}

func initEnv(path *string) (err error) {
	if loadedEnv {
		return nil
	}
	p := ""
	if path == nil {
		p = "../.env"
	} else {
		p = *path
	}
	loadedEnv = true
	return godotenv.Load(p)
}
