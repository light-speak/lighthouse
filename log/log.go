package log

import (
	"os"

	"github.com/rs/zerolog"
)

var Log *zerolog.Logger

func Trace() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Trace()
}

func Debug() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Debug()
}

func Info() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Info()
}

func Warn() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Warn()
}

func Error() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Error()
}

func Fatal() *zerolog.Event {
	if Log == nil {
		NewLogger()
	}
	return Log.Fatal()
}

func NewLogger() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	Log = &logger
}
