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

const minWidth = 20 // 设置默认最小宽度

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
}

// LoggerConfig 是用于初始化 Logger 的配置结构体
type LoggerConfig struct {
	LogDir      string // 日志存放的根目录
	ConsoleOnly bool   // 是否只输出到控制台
	FileOnly    bool   // 是否只输出到文件
	Level       int    // 日志输出的最低级别
}

// NewLogger 创建一个新的 Logger 实例，并初始化日志文件和目录
func NewLogger(config LoggerConfig) (*Logger, error) {
	var file *os.File
	var err error

	// 确定日志文件夹路径并创建
	logDir := filepath.Join(config.LogDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建日志目录: %v", err)
	}

	// 如果不是 ConsoleOnly 模式，初始化日志文件
	if !config.ConsoleOnly {
		file, err = openDailyFile(logDir)
		if err != nil {
			return nil, err
		}
	}

	// 返回初始化的 Logger 实例
	return &Logger{
		LogFile:     file,
		ConsoleOnly: config.ConsoleOnly,
		FileOnly:    config.FileOnly,
		logDir:      logDir,
		Level:       config.Level,
	}, nil
}

// openDailyFile 打开或创建今天的日志文件，并清理超过15天的旧日志文件
func openDailyFile(logDir string) (*os.File, error) {
	today := time.Now().Format("2006-01-02")
	dailyLogFilePath := filepath.Join(logDir, fmt.Sprintf("log-%s.log", today))

	// 打开或创建今天的日志文件
	file, err := os.OpenFile(dailyLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 清理超过15天的旧日志文件
	err = cleanupOldLogs(logDir, 15)
	if err != nil {
		fmt.Fprintf(os.Stderr, "清理旧日志文件时出错: %v\n", err)
	}

	return file, nil
}

// cleanupOldLogs 清理超过指定数量的旧日志文件
func cleanupOldLogs(logDir string, maxDays int) error {
	files, err := filepath.Glob(filepath.Join(logDir, "log-*.log"))
	if err != nil {
		return fmt.Errorf("无法列出日志文件: %v", err)
	}

	// 如果文件数量少于或等于 maxDays，则无需删除任何文件
	if len(files) <= maxDays {
		return nil
	}

	// 按文件名（日期）排序
	sort.Strings(files)

	// 删除最早的文件，保留最近的 maxDays 个文件
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
	_, file, line, ok := runtime.Caller(3)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown"
}

// formatLogMessage 格式化日志信息，并处理多行显示
func formatLogMessage(level string, message string, callerInfo string) (string, string, string) {
	lines := strings.Split(message, "\n")
	maxWidth := minWidth

	// 计算最大行宽
	for _, line := range lines {
		lineWidth := calcDisplayWidth(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	// 确保maxWidth不小于调用信息的宽度
	callerInfoWidth := calcDisplayWidth(callerInfo) + 12
	if callerInfoWidth > maxWidth {
		maxWidth = callerInfoWidth
	}

	// 设置标题宽度为最长行的宽度
	levelLine := fmt.Sprintf("╓%s", centerText(fmt.Sprintf(" %s ", level), maxWidth+2))
	footerLine := fmt.Sprintf("╙%s", centerText("LIGHTHOUSE", maxWidth+2))

	// 添加调用信息行
	callerLine := fmt.Sprintf("║ Called from %s%s ", callerInfo, strings.Repeat(" ", maxWidth-calcDisplayWidth(callerInfo)-12))

	// 格式化每一行
	var formattedLines []string
	for _, line := range lines {
		padding := maxWidth - calcDisplayWidth(line)
		formattedLines = append(formattedLines, fmt.Sprintf("║ %s%s ", line, strings.Repeat(" ", padding)))
	}

	return levelLine, fmt.Sprintf("%s\n%s", callerLine, strings.Join(formattedLines, "\n")), footerLine
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
		fileMessage := fmt.Sprintf("[%s] [%s] %s - %s", timestamp, levelName, callerInfo, message)
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
