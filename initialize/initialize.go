package initialize

import (
	"fmt"
	"github.com/light-speak/lighthouse/log"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func Run() error {
	if err := createFiles(); err != nil {
		return err
	}
	return nil
}

func createFiles() error {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 获取当前目录的名称
	dirName := filepath.Base(dir)

	// 构造 .env 文件的路径
	envFilePath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envFilePath); err == nil {
		log.Warn(".env file already exists")
		return nil
	}

	// 创建 .env 文件
	file, err := os.Create(envFilePath)
	if err != nil {
		log.Error("Failed to create .env file: %s", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("Error closing file: %s", err)
		}
	}(file)

	// 写入默认的 .env 内容
	_, err = file.WriteString(fmt.Sprintf("# Default .env file"+
		"\n\nDB_HOST=127.0.0.1\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=\nDB_NAME=%s\nDB_LOG_LEVEL=INFO"+
		"\n\nLOG_LEVEL=INFO\nLOG_PATH=storage"+
		"\n\nPORT=4001", dirName))
	if err != nil {
		log.Error("Failed to write to .env file: %s", err)
		return err
	}

	log.Info(".env file created successfully")

	// 复制 ../tpl/.gitignore 到当前目录并重命名为 .gitignore
	err = copyGitignoreFile(dir)
	if err != nil {
		log.Error("Failed to copy .gitignore file: %s", err)
		return err
	}

	log.Info(".gitignore file copied successfully")
	return nil
}

// getLibraryPath 获取当前库的根目录路径
func getLibraryPath() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// b 是当前文件的路径, 返回 lighthouse 目录
	return filepath.Dir(filepath.Dir(b)), nil
}

// copyGitignoreFile 复制库目录下的 .gitignore 文件到工作目录
func copyGitignoreFile(destDir string) error {
	// 获取 lighthouse 包的根目录

	// 获取库的根目录路径
	libDir, err := getLibraryPath()
	if err != nil {
		return err
	}

	// 构造库中 tpl/.gitignore 文件的路径
	srcFilePath := filepath.Join(libDir, "tpl", ".gitignore")

	// 构造目标文件路径
	destFilePath := filepath.Join(destDir, ".gitignore")

	// 检查目标文件是否已经存在
	if _, err := os.Stat(destFilePath); err == nil {
		log.Warn(".gitignore file already exists, skipping copy")
		return nil
	}

	// 打开源文件
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		log.Error("Failed to open source file: %s", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(destFilePath)
	if err != nil {
		log.Error("Failed to create destination file: %s", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// 执行文件复制
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Error("Failed to copy file: %s", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// 确保文件内容成功写入磁盘
	err = destFile.Sync()
	if err != nil {
		log.Error("Failed to sync destination file: %s", err)
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
