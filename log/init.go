package log

import (
	"fmt"
	"github.com/light-speak/lighthouse/env"
	"os"
	"sync"
)

var (
	logger *Logger
	once   sync.Once
)

const (
	defaultLogLevel = "DEBUG"
	defaultLogPath  = "storage"
)

func initLogger() {
	once.Do(func() {
		logLevelStr := env.Getenv("LOG_LEVEL", defaultLogLevel)
		logDir := env.Getenv("LOG_PATH", defaultLogPath)

		logLevel := parseLogLevel(logLevelStr)

		config := LoggerConfig{
			LogDir:      logDir,
			ConsoleOnly: false,
			FileOnly:    false,
			Level:       logLevel,
		}

		var err error
		logger, err = NewLogger(config)
		if err != nil {
			fmt.Printf("初始化日志器失败: %v\n", err)
			os.Exit(1)
		}
	})
}

func parseLogLevel(levelStr string) int {
	switch levelStr {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelDebug
	}
}

func Debug(format string, args ...interface{}) {
	initLogger()
	callerInfo := getCallerInfo()
	logger.debug(format, callerInfo, args...)
}

func Info(format string, args ...interface{}) {
	initLogger()
	callerInfo := getCallerInfo()
	logger.info(format, callerInfo, args...)
}

func Warn(format string, args ...interface{}) {
	initLogger()
	callerInfo := getCallerInfo()
	logger.warn(format, callerInfo, args...)
}

func Error(format string, args ...interface{}) {
	initLogger()
	callerInfo := getCallerInfo()
	logger.error(format, callerInfo, args...)
}
