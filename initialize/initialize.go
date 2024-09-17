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
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	dirName := filepath.Base(dir)

	if err := createEnvFile(dir, dirName); err != nil {
		return err
	}

	if err := copyGitignoreFile(dir); err != nil {
		return err
	}

	return nil
}

func createEnvFile(dir, dirName string) error {
	envFilePath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envFilePath); err == nil {
		log.Warn(".env 文件已存在")
		return nil
	}

	file, err := os.Create(envFilePath)
	if err != nil {
		return fmt.Errorf("创建 .env 文件失败: %w", err)
	}
	defer file.Close()

	envContent := fmt.Sprintf(`# 默认 .env 文件

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=%s
DB_LOG_LEVEL=INFO

LOG_LEVEL=INFO
LOG_PATH=storage

PORT=4001`, dirName)

	if _, err := file.WriteString(envContent); err != nil {
		return fmt.Errorf("写入 .env 文件失败: %w", err)
	}

	log.Info(".env 文件创建成功")
	return nil
}

func getLibraryPath() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}
	return filepath.Dir(filepath.Dir(b)), nil
}

func copyGitignoreFile(destDir string) error {
	libDir, err := getLibraryPath()
	if err != nil {
		return err
	}

	srcFilePath := filepath.Join(libDir, "tpl", ".gitignore")
	destFilePath := filepath.Join(destDir, ".gitignore")

	if _, err := os.Stat(destFilePath); err == nil {
		log.Warn(".gitignore 文件已存在，跳过复制")
		return nil
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("同步目标文件失败: %w", err)
	}

	log.Info(".gitignore 文件复制成功")
	return nil
}
