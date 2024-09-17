package log

import (
	"fmt"
	"github.com/light-speak/lighthouse/env"
	"os"
	"sync"
)

var logger *Logger
var once sync.Once

func initLogger() {
	once.Do(func() {
		// 从环境变量获取日志级别和路径
		logLevelStr := env.Getenv("LOG_LEVEL", "DEBUG")
		logDir := env.Getenv("LOG_PATH", "storage")

		logLevel := LevelDebug
		switch logLevelStr {
		case "DEBUG":
			logLevel = LevelDebug
			break
		case "INFO":
			logLevel = LevelInfo
			break
		case "WARN":
			logLevel = LevelWarn
			break
		case "ERROR":
			logLevel = LevelError
			break
		}

		// 初始化 Logger 实例
		config := LoggerConfig{
			LogDir:      logDir,   // 日志目录
			ConsoleOnly: false,    // 是否只输出到控制台
			FileOnly:    false,    // 是否只输出到文件
			Level:       logLevel, // 日志级别
		}

		var err error
		logger, err = NewLogger(config)
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
	})
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
