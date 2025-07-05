package databases

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm/logger"
)

var databaseConfig *DatabaseConfig

type DatabaseConfig struct {
	Hosts    []string
	Port     string
	User     string
	Password string
	Name     string
	LogLevel logger.LogLevel

	// 1.0 版本，支持多数据库，同时兼容原有数据库配置
	EnableSlave bool

	Main  *DatabaseConfig
	Slave *DatabaseConfig
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func parseHosts(hostStr string) []string {
	if hostStr == "" {
		return []string{"localhost"}
	}
	return strings.Split(hostStr, ",")
}

func init() {
	databaseConfig = &DatabaseConfig{
		Hosts:    []string{"localhost"},
		Port:     "3306",
		User:     "root",
		Password: "",
		Name:     "example",
		LogLevel: logger.Info,
		Main: &DatabaseConfig{
			Hosts:    []string{"localhost"},
			Port:     "3306",
			User:     "root",
			Password: "",
			Name:     "example",
		},
		Slave: &DatabaseConfig{
			Hosts:    []string{"localhost"},
			Port:     "3306",
			User:     "root",
			Password: "",
			Name:     "example",
		},
	}

	if curPath, err := os.Getwd(); err == nil {
		_ = godotenv.Load(filepath.Join(curPath, ".env"))
	}

	databaseConfig.EnableSlave = utils.GetEnv("DB_ENABLE_SLAVE", "false") == "true"

	if databaseConfig.EnableSlave {
		databaseConfig.Main.Hosts = parseHosts(utils.GetEnv("DB_MAIN_HOST", strings.Join(databaseConfig.Main.Hosts, ",")))
		databaseConfig.Main.Port = utils.GetEnv("DB_MAIN_PORT", databaseConfig.Main.Port)
		databaseConfig.Main.User = utils.GetEnv("DB_MAIN_USER", databaseConfig.Main.User)
		databaseConfig.Main.Password = utils.GetEnv("DB_MAIN_PASSWORD", databaseConfig.Main.Password)
		databaseConfig.Main.Name = utils.GetEnv("DB_NAME", databaseConfig.Main.Name)

		databaseConfig.Slave.Hosts = parseHosts(utils.GetEnv("DB_SLAVE_HOST", strings.Join(databaseConfig.Slave.Hosts, ",")))
		databaseConfig.Slave.Port = utils.GetEnv("DB_SLAVE_PORT", databaseConfig.Slave.Port)
		databaseConfig.Slave.User = utils.GetEnv("DB_SLAVE_USER", databaseConfig.Slave.User)
		databaseConfig.Slave.Password = utils.GetEnv("DB_SLAVE_PASSWORD", databaseConfig.Slave.Password)
		databaseConfig.Slave.Name = utils.GetEnv("DB_NAME", databaseConfig.Slave.Name)

	} else {
		databaseConfig.Main.Hosts = parseHosts(utils.GetEnv("DB_HOST", strings.Join(databaseConfig.Hosts, ",")))
		databaseConfig.Main.Port = utils.GetEnv("DB_PORT", databaseConfig.Port)
		databaseConfig.Main.User = utils.GetEnv("DB_USER", databaseConfig.User)
		databaseConfig.Main.Password = utils.GetEnv("DB_PASSWORD", databaseConfig.Password)
		databaseConfig.Main.Name = utils.GetEnv("DB_NAME", databaseConfig.Name)

		// 兼容一下单服务但是填的是Main的情况
		mainHost := utils.GetEnv("DB_MAIN_HOST", "-")
		if mainHost != "-" {
			databaseConfig.Main.Hosts = parseHosts(mainHost)
		}
		mainPort := utils.GetEnv("DB_MAIN_PORT", "-")
		if mainPort != "-" {
			databaseConfig.Main.Port = mainPort
		}
		mainUser := utils.GetEnv("DB_MAIN_USER", "-")
		if mainUser != "-" {
			databaseConfig.Main.User = mainUser
		}
		mainPassword := utils.GetEnv("DB_MAIN_PASSWORD", "-")
		if mainPassword != "-" {
			databaseConfig.Main.Password = mainPassword
		}
	}

	switch LogLevel(utils.GetEnv("DB_LOG_LEVEL", string(LogLevelInfo))) {
	case LogLevelDebug:
		databaseConfig.LogLevel = logger.Info
	case LogLevelInfo:
		databaseConfig.LogLevel = logger.Info
	case LogLevelWarn:
		databaseConfig.LogLevel = logger.Warn
	case LogLevelError:
		databaseConfig.LogLevel = logger.Error
	default:
		databaseConfig.LogLevel = logger.Info
	}

}
