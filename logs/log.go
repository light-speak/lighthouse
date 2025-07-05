package logs

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/light-speak/lighthouse/utils"
	"github.com/rs/zerolog"
)

var (
	Log            *zerolog.Logger
	currentOutputs []io.Writer
	currentLogFile *os.File
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
	*outputs = append(*outputs, fileWriter)
	return nil
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
