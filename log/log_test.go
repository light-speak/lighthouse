package log

import (
	"testing"

	"github.com/rs/zerolog/log"
)

func TestInit(t *testing.T) {
	log.Info().Msg("log init")
}
