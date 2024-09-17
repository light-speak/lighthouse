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
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}

	libDir, err := getLibraryPath()
	if err != nil {
		return fmt.Errorf("获取库路径失败: %w", err)
	}

	itemsToCopy := []struct {
		src  string
		dest string
	}{
		{filepath.Join(libDir, "../tpl", "graph"), filepath.Join(currentDir, "graph")},
		{filepath.Join(libDir, "../tpl", "gqlgen.yml"), filepath.Join(currentDir, "gqlgen.yml")},
	}

	for _, item := range itemsToCopy {
		if err := copyItem(item.src, item.dest); err != nil {
			return err
		}
	}

	return nil
}

func getLibraryPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}
	return filepath.Dir(filepath.Dir(currentFilePath)), nil
}

func copyItem(src, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		log.Warn("%s 已存在，跳过复制", dest)
		return nil
	}

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	if info.IsDir() {
		return copyDir(src, dest)
	}
	return copyFile(src, dest)
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	if err = destFile.Sync(); err != nil {
		return fmt.Errorf("同步目标文件失败: %w", err)
	}

	log.Info("已复制文件: %s -> %s", src, dest)
	return nil
}

func copyDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
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

	log.Info("已复制目录: %s -> %s", src, dest)
	return nil
}
