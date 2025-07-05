package logs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level      zerolog.Level
	TimeFormat string
	Caller     bool
	Console    bool
	File       bool
	Pretty     bool
	FilePath   string
}

var loggerConfig *LoggerConfig

func InitLogger() error {
	loggerConfig = &LoggerConfig{
		Level:      zerolog.InfoLevel,
		TimeFormat: "2006-01-02 15:04:05",
		Caller:     false,
		Console:    true,
		File:       false,
		FilePath:   "logs/logs.log",
		Pretty:     true,
	}

	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}
	loggerConfig.Level = getLogLevel(utils.GetEnv("LOG_LEVEL", "info"))
	loggerConfig.TimeFormat = utils.GetEnv("LOG_TIME_FORMAT", loggerConfig.TimeFormat)
	loggerConfig.Caller = utils.GetEnvBool("LOG_CALLER", loggerConfig.Caller)
	loggerConfig.Console = utils.GetEnvBool("LOG_CONSOLE", loggerConfig.Console)
	loggerConfig.File = utils.GetEnvBool("LOG_FILE", loggerConfig.File)
	loggerConfig.FilePath = utils.GetEnv("LOG_FILE_PATH", loggerConfig.FilePath)
	loggerConfig.Pretty = utils.GetEnvBool("LOG_PRETTY", loggerConfig.Pretty)

	currentOutputs = nil // Reset outputs on init
	return setupLogger()
}

func getLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
