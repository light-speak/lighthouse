package env

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/plugins/elasticsearch"
	"github.com/light-speak/lighthouse/plugins/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type LoggerOutput interface {
	Write(p []byte) (n int, err error)
}

// InitLogger Initialize the logger
func InitLogger() {
	var output LoggerOutput

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(LighthouseConfig.Logger.Level)
	if LighthouseConfig.Logger.Stack {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	switch LighthouseConfig.Logger.Driver {
	case Stdout:
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	case File:
		fileOutput, err := file.NewFileOutput(LighthouseConfig.Logger.Path)
		if err != nil {
			output = os.Stdout
			fmt.Printf("Failed to initialize file output: %v", err)
		} else {
			output = zerolog.MultiLevelWriter(fileOutput, os.Stdout)
		}
	case Elasticsearch:
		esOutput, err := elasticsearch.NewElasticsearchOutput()
		if err != nil {
			fmt.Printf("Failed to initialize elasticsearch output: %v", err)
			output = zerolog.MultiLevelWriter(os.Stdout, esOutput)
		} else {
			output = esOutput
		}
	default:
		// Default to stdout if an unknown driver is specified
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	cusLog := zerolog.New(output).With().Timestamp().Logger()
	log.Log = &cusLog
	log.Log.Info().Msgf("Logger initialized, level: %s", LighthouseConfig.Logger.Level.String())
}
