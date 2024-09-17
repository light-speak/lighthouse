package initialize

import (
	"fmt"

	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/light-speak/lighthouse/log"
)

func Run() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// 获取库的根目录路径
	libDir, err := getLibraryPath()
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	// 需要复制的文件和文件夹
	itemsToCopy := []struct {
		src  string
		dest string
	}{
		{src: filepath.Join(libDir, "../tpl", "graph"), dest: filepath.Join(currentDir, "graph")},
		{src: filepath.Join(libDir, "../tpl", "gqlgen.yml"), dest: filepath.Join(currentDir, "gqlgen.yml")},
	}

	// 复制每个文件/文件夹
	for _, item := range itemsToCopy {
		if err := copyItem(item.src, item.dest); err != nil {
			return err
		}
	}

	return nil
}

// getLibraryPath 获取当前库的根目录路径
func getLibraryPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// b 是当前文件的路径, 返回 lighthouse 目录
	return filepath.Dir(filepath.Dir(currentFilePath)), nil
}

// copyItem 复制文件或文件夹，如果已存在则跳过
func copyItem(src, dest string) error {
	// 检查目标文件/文件夹是否已经存在
	if _, err := os.Stat(dest); err == nil {
		log.Warn("%s already exists, skipping copy", dest)
		return nil
	}

	// 检查是文件还是文件夹
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if info.IsDir() {
		return copyDir(src, dest)
	}
	return copyFile(src, dest)
}

// copyFile 复制文件
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Error("Failed to open source file: %s", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		log.Error("Failed to create destination file: %s", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Error("Failed to copy file: %s", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		log.Error("Failed to sync destination file: %s", err)
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	log.Info("Copied file from %s to %s", src, dest)
	return nil
}

// copyDir 递归复制文件夹
func copyDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		log.Error("Failed to read directory: %s", err)
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		log.Error("Failed to create directory: %s", err)
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if filepath.Ext(entry.Name()) == ".gotpl" {
				destPath = filepath.Join(dest, strings.TrimSuffix(entry.Name(), ".gotpl")+".go")
			}

			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	log.Info("Copied directory from %s to %s", src, dest)
	return nil
}
