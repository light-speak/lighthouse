package lighterr

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/utils"
)

type ErrorConfig struct {
	Env Env
}

type Env string

const (
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)

var config *ErrorConfig

func init() {
	config = &ErrorConfig{
		Env: EnvDevelopment,
	}

	if cp, err := os.Getwd(); err == nil {
		_ = godotenv.Load(filepath.Join(cp, ".env"))
	}

	config.Env = Env(utils.GetEnv("APP_ENV", string(config.Env)))
	logs.Debug().Msgf("env: %s", config.Env)
}
