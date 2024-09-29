package env

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var LighthouseConfig *Config

type AppEnvironment string

const (
	Development AppEnvironment = "development"
	Production  AppEnvironment = "production"
	Staging     AppEnvironment = "staging"
)

type AppMode string

const (
	Single  AppMode = "single"
	Cluster AppMode = "cluster"
)

type DatabaseDriver string

const (
	Postgres DatabaseDriver = "postgres"
	MySQL    DatabaseDriver = "mysql"
)

type DatabaseOrm string

const (
	Gorm DatabaseOrm = "gorm"
	Ent  DatabaseOrm = "ent"
)

type AuthDriver string

const (
	Jwt    AuthDriver = "jwt"
	Router AuthDriver = "router"
)

type LoggerLevel string

type Config struct {
	App struct {
		Name        string
		Environment AppEnvironment
		Mode        AppMode
	}
	Server struct {
		Port string
	}
	Api struct {
		Restful bool
		Prefix  string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		Driver   DatabaseDriver
		Orm      DatabaseOrm
	}
	Auth struct {
		Driver AuthDriver
	}
	Logger struct {
		Level    zerolog.Level
		Path     string
		Beautify bool
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		Db       int
	}
}

// InitConfig Initialize the configuration
func init() {
	curPath, err := os.Getwd()
	if err != nil {
		return
	}
	godotenv.Load(filepath.Join(curPath, ".env"))
	LighthouseConfig = &Config{
		App: struct {
			Name        string
			Environment AppEnvironment
			Mode        AppMode
		}{
			Name:        GetEnv("APP_NAME", "MyApp"),
			Environment: AppEnvironment(GetEnv("APP_ENVIRONMENT", "development")),
			Mode:        AppMode(GetEnv("APP_MODE", "single")),
		},
		Server: struct {
			Port string
		}{
			Port: GetEnv("SERVER_PORT", "8080"),
		},
		Api: struct {
			Restful bool
			Prefix  string
		}{
			Restful: GetEnvBool("API_RESTFUL", false),
			Prefix:  GetEnv("API_PREFIX", "/api"),
		},
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
			Driver   DatabaseDriver
			Orm      DatabaseOrm
		}{
			Host:     GetEnv("DB_HOST", "127.0.0.1"),
			Port:     GetEnv("DB_PORT", "3306"),
			User:     GetEnv("DB_USER", "root"),
			Password: GetEnv("DB_PASSWORD", ""),
			Name:     GetEnv("DB_NAME", "example"),
			Driver:   DatabaseDriver(GetEnv("DB_DRIVER", "mysql")),
			Orm:      DatabaseOrm(GetEnv("DB_ORM", "gorm")),
		},
		Auth: struct {
			Driver AuthDriver
		}{
			Driver: AuthDriver(GetEnv("AUTH_DRIVER", "jwt")),
		},
		Logger: struct {
			Level    zerolog.Level
			Path     string
			Beautify bool
		}{
			Level:    zerolog.Level(GetEnvInt("LOGGER_LEVEL", 0)),
			Path:     GetEnv("LOGGER_PATH", "logs/app.log"),
			Beautify: GetEnvBool("LOGGER_BEAUTIFY", true),
		},
		Redis: struct {
			Host     string
			Port     string
			Password string
			Db       int
		}{
			Host:     GetEnv("REDIS_HOST", "localhost"),
			Port:     GetEnv("REDIS_PORT", "6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			Db:       GetEnvInt("REDIS_DB", 0),
		},
	}
}
