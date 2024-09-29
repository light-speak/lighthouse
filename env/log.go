package env

import (
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
		output = os.Stdout
	case File:
		fileOutput, err := file.NewFileOutput(LighthouseConfig.Logger.Path)
		if err != nil {
			log.Error().Msgf("Failed to initialize file output: %v", err)
			output = os.Stdout
		} else {
			output = fileOutput
		}
	case Elasticsearch:
		esOutput, err := elasticsearch.NewElasticsearchOutput()
		if err != nil {
			log.Error().Msgf("Failed to initialize elasticsearch output: %v", err)
			output = os.Stdout
		} else {
			output = esOutput
		}
	default:
		// Default to stdout if an unknown driver is specified
		output = os.Stdout
	}

	log.Log = zerolog.New(output).With().Timestamp().Logger()
}
