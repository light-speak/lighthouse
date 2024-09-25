package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	minWidth        = 20
	maxLogRetention = 15
	logFilePattern  = "log-*.log"
)

// 定义日志级别常量
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger 是自定义日志结构体，包含日志文件和相关配置
type Logger struct {
	mu          sync.Mutex // 互斥锁，确保日志写入时的并发安全
	LogFile     *os.File   // 日志文件句柄
	UseColor    bool       // 是否使用彩色输出（未使用）
	ConsoleOnly bool       // 是否只输出到控制台
	FileOnly    bool       // 是否只输出到文件
	logDir      string     // 日志文件存放目录
	Level       int        // 日志输出的最低级别
	ShowCaller  bool       // 是否显示调用者信息
}

// LoggerConfig 是用于初始化 Logger 的配置结构体
type LoggerConfig struct {
	LogDir      string // 日志存放的根目录
	ConsoleOnly bool   // 是否只输出到控制台
	FileOnly    bool   // 是否只输出到文件
	Level       int    // 日志输出的最低级别
	ShowCaller  bool   // 是否显示调用者信息
}

// NewLogger 创建一个新的 Logger 实例，并初始化日志文件和目录
func NewLogger(config LoggerConfig) (*Logger, error) {
	// 获取项目根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("无法找到项目根目录: %v", err)
	}

	// 在项目根目录下创建日志目录
	logDir := filepath.Join(projectRoot, config.LogDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建日志目录: %v", err)
	}

	var file *os.File
	if !config.ConsoleOnly {
		file, err = openDailyFile(logDir)
		if err != nil {
			return nil, err
		}
	}

	return &Logger{
		LogFile:     file,
		ConsoleOnly: config.ConsoleOnly,
		FileOnly:    config.FileOnly,
		logDir:      logDir,
		Level:       config.Level,
		ShowCaller:  config.ShowCaller,
	}, nil
}

// findProjectRoot 查找项目根目录
func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 首先检查当前目录是否存在 go.mod 文件
	if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
		return currentDir, nil
	}

	// 如果没有 go.mod 文件，尝试获取可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 获取可执行文件所在的目录
	execDir := filepath.Dir(execPath)

	// 如果可执行文件目录存在 go.mod，则返回该目录
	if _, err := os.Stat(filepath.Join(execDir, "go.mod")); err == nil {
		return execDir, nil
	}

	// 如果都没有找到 go.mod，则返回可执行文件所在的目录作为项目根目录
	return execDir, nil
}

// openDailyFile 打开或创建今天的日志文件，并清理超过15天的旧日志文件
func openDailyFile(logDir string) (*os.File, error) {
	today := time.Now().Format("2006-01-02")
	dailyLogFilePath := filepath.Join(logDir, fmt.Sprintf("log-%s.log", today))

	file, err := os.OpenFile(dailyLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	if err := cleanupOldLogs(logDir, maxLogRetention); err != nil {
		fmt.Fprintf(os.Stderr, "清理旧日志文件时出错: %v\n", err)
	}

	return file, nil
}

// cleanupOldLogs 清理超过指定数量的旧日志文件
func cleanupOldLogs(logDir string, maxDays int) error {
	files, err := filepath.Glob(filepath.Join(logDir, logFilePattern))
	if err != nil {
		return fmt.Errorf("无法列出日志文件: %v", err)
	}

	if len(files) <= maxDays {
		return nil
	}

	sort.Strings(files)

	for _, file := range files[:len(files)-maxDays] {
		if err := os.Remove(file); err != nil {
			fmt.Fprintf(os.Stderr, "无法删除旧日志文件: %v\n", err)
		}
	}

	return nil
}

// calcDisplayWidth 计算字符串的实际显示宽度，中文字符占2个单位宽度
func calcDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if utf8.RuneLen(r) == 3 {
			width += 2 // 中文字符宽度为2
		} else {
			width += 1 // 英文字符宽度为1
		}
	}
	return width
}

// centerText 居中显示文本，宽度不足的地方填充 "═"
func centerText(text string, width int) string {
	textWidth := calcDisplayWidth(text)
	padding := (width - textWidth) / 2
	if padding > 0 {
		return fmt.Sprintf("%s%s%s", strings.Repeat("═", padding), text, strings.Repeat("═", padding))
	}
	return fmt.Sprintf("%s%s", text, strings.Repeat("═", width-textWidth))
}

// getCallerInfo 获取调用者的信息，包含文件名和行号
func getCallerInfo() string {
	if !(*logger).ShowCaller {
		return ""
	}
	for i := 1; i < 10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.Contains(file, "lighthouse/log") {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return "未知"
}

// formatLogMessage 格式化日志信息，并处理多行显示
func formatLogMessage(level, message, callerInfo string) (string, string, string) {
	lines := strings.Split(message, "\n")
	maxWidth := minWidth

	for _, line := range lines {
		lineWidth := calcDisplayWidth(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	if callerInfo != "" {
		callerInfoWidth := calcDisplayWidth(callerInfo) + 12
		if callerInfoWidth > maxWidth {
			maxWidth = callerInfoWidth
		}
	}

	color := getColorForLevel(level)
	resetColor := "\033[0m"

	levelLine := fmt.Sprintf("%s╓%s%s", color, centerText(fmt.Sprintf(" %s ", level), maxWidth+2), resetColor)
	footerLine := fmt.Sprintf("%s╙%s%s", color, centerText("LIGHTHOUSE", maxWidth+2), resetColor)

	var formattedLines []string
	if callerInfo != "" {
		callerLine := fmt.Sprintf("%s║ Called from %s%s %s", color, callerInfo, strings.Repeat(" ", maxWidth-calcDisplayWidth(callerInfo)-12), resetColor)
		formattedLines = append(formattedLines, callerLine)
	}

	for _, line := range lines {
		padding := maxWidth - calcDisplayWidth(line)
		formattedLines = append(formattedLines, fmt.Sprintf("%s║ %s%s %s", color, line, strings.Repeat(" ", padding), resetColor))
	}

	return levelLine, strings.Join(formattedLines, "\n"), footerLine
}

func getColorForLevel(level string) string {
	switch level {
	case "DEBUG ":
		return "\033[36m"
	case " INFO ":
		return "\033[32m"
	case " WARN ":
		return "\033[33m"
	case "ERROR ":
		return "\033[31m"
	default:
		return "\033[0m"
	}
}

// Log 记录日志信息，带有日志级别限制功能，并支持格式化字符串
func (l *Logger) Log(level int, levelName, format string, callerInfo string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 如果日志级别低于设定的最低级别，不记录日志
	if level < l.Level {
		return
	}

	// 格式化消息
	message := fmt.Sprintf(format, args...)

	// 格式化日志消息
	levelLine, formattedMessage, footerLine := formatLogMessage(levelName, message, callerInfo)

	// 输出到文件
	if l.LogFile != nil && !l.ConsoleOnly {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		var fileMessage string
		if callerInfo != "" {
			fileMessage = fmt.Sprintf("[%s] [%s] %s - %s", timestamp, levelName, callerInfo, message)
		} else {
			fileMessage = fmt.Sprintf("[%s] [%s] %s", timestamp, levelName, message)
		}
		fmt.Fprintln(l.LogFile, fileMessage)
	}

	// 输出到控制台
	if !l.FileOnly {
		fmt.Println(levelLine)
		fmt.Println(formattedMessage)
		fmt.Println(footerLine)
	}
}

// Debug 记录调试级别的日志，支持格式化字符串
func (l *Logger) debug(format string, callerInfo string, args ...interface{}) {
	l.Log(LevelDebug, "DEBUG ", format, callerInfo, args...)
}

// Info 记录信息级别的日志，支持格式化字符串
func (l *Logger) info(format string, callerInfo string, args ...interface{}) {
	l.Log(LevelInfo, " INFO ", format, callerInfo, args...)
}

// Error 记录错误级别的日志，支持格式化字符串
func (l *Logger) error(format string, callerInfo string, args ...interface{}) {
	l.Log(LevelError, "ERROR ", format, callerInfo, args...)
}

// Warn 记录警告级别的日志，支持格式化字符串
func (l *Logger) warn(format string, callerInfo string, args ...interface{}) {
	l.Log(LevelWarn, " WARN ", format, callerInfo, args...)
}
