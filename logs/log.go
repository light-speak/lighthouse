package logs

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/light-speak/lighthouse/utils"
	"github.com/rs/zerolog"
)

var (
	Log            *zerolog.Logger
	currentOutputs []io.Writer
	currentLogFile *os.File
	currentLogDate string       // 当前日志文件的日期
	logMutex       sync.RWMutex // 保护日志切换的并发安全
	rotateOnce     sync.Once    // 确保日志轮转 goroutine 只启动一次
)

func SetOutput(out ...io.Writer) error {
	if len(out) > 0 {
		currentOutputs = append(currentOutputs, out...)
	}
	return setupLogger()
}

func setupLogger() error {
	var outputs []io.Writer

	// Console output
	if loggerConfig.Console {
		var consoleWriter io.Writer
		if loggerConfig.Pretty {
			consoleWriter = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: loggerConfig.TimeFormat,
			}
			outputs = append(outputs, consoleWriter)
		} else {
			consoleWriter = os.Stderr
			outputs = append(outputs, consoleWriter)
		}
	}

	// File output
	if loggerConfig.File {
		if err := setupFileOutput(&outputs); err != nil {
			return err
		}
	}

	// Add custom outputs
	if len(currentOutputs) > 0 {
		outputs = append(outputs, currentOutputs...)
	}

	// Create multi-writer
	var writer io.Writer
	if len(outputs) > 1 {
		writer = io.MultiWriter(outputs...)
	} else if len(outputs) == 1 {
		writer = outputs[0]
	} else {
		writer = os.Stdout
	}

	zerolog.SetGlobalLevel(loggerConfig.Level)

	logContext := zerolog.New(writer).With().Timestamp()
	if loggerConfig.Caller {
		logContext = logContext.Caller()
	}
	logger := logContext.Logger()

	Log = &logger
	return nil
}

func setupFileOutput(outputs *[]io.Writer) error {
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Dir(loggerConfig.FilePath)
	fileName := filepath.Base(loggerConfig.FilePath)
	ext := filepath.Ext(fileName)
	name := fileName[:len(fileName)-len(ext)]
	dailyLogPath := filepath.Join(logDir, name+"-"+today+ext)

	if err := utils.MkdirAll(logDir); err != nil {
		return err
	}

	// Close existing log file if it exists
	if currentLogFile != nil {
		currentLogFile.Close()
	}

	fileWriter, err := os.OpenFile(dailyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	currentLogFile = fileWriter
	currentLogDate = today
	*outputs = append(*outputs, fileWriter)

	// 启动日志轮转检查 goroutine（只启动一次）
	rotateOnce.Do(func() {
		go startLogRotation()
	})

	return nil
}

// startLogRotation 定时检查是否需要切换日志文件
func startLogRotation() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if !loggerConfig.File {
			continue
		}

		today := time.Now().Format("2006-01-02")
		logMutex.RLock()
		needRotate := currentLogDate != today
		logMutex.RUnlock()

		if needRotate {
			logMutex.Lock()
			// 双重检查
			if currentLogDate != today {
				rotateLogFile()
			}
			logMutex.Unlock()
		}
	}
}

// rotateLogFile 切换到新的日志文件
func rotateLogFile() {
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Dir(loggerConfig.FilePath)
	fileName := filepath.Base(loggerConfig.FilePath)
	ext := filepath.Ext(fileName)
	name := fileName[:len(fileName)-len(ext)]
	dailyLogPath := filepath.Join(logDir, name+"-"+today+ext)

	// 关闭旧文件
	if currentLogFile != nil {
		currentLogFile.Close()
	}

	// 打开新文件
	fileWriter, err := os.OpenFile(dailyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 出错时继续使用旧配置，记录到 stderr
		os.Stderr.WriteString("日志轮转失败: " + err.Error() + "\n")
		return
	}

	currentLogFile = fileWriter
	currentLogDate = today

	// 重新设置 logger
	setupLogger()
}

func Trace() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Trace()
}

func Debug() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Debug()
}

func Info() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Info()
}

func Warn() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Warn()
}

func Error() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Error()
}

func Fatal() *zerolog.Event {
	if Log == nil {
		if err := InitLogger(); err != nil {
			panic(err)
		}
	}
	return Log.Fatal()
}
