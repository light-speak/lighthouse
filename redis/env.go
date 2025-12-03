package redis

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
)

// # Redis settings
// REDIS_HOST=localhost
// REDIS_PORT=6379
// REDIS_PASSWORD=
// REDIS_DB=0
type Config struct {
	Enable   bool
	Host     string
	Port     string
	Password string
	DB       int

	// 连接池配置
	PoolSize     int // 连接池大小
	MinIdleConns int // 最小空闲连接数
}

var LightRedisConfig *Config

func init() {
	LightRedisConfig = &Config{
		Enable:   false,
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	LightRedisConfig.Enable = utils.GetEnvBool("REDIS_ENABLE", false)
	LightRedisConfig.Host = utils.GetEnv("REDIS_HOST", "localhost")
	LightRedisConfig.Port = utils.GetEnv("REDIS_PORT", "6379")
	LightRedisConfig.Password = utils.GetEnv("REDIS_PASSWORD", "")
	LightRedisConfig.DB = utils.GetEnvInt("REDIS_DB", 0)
	LightRedisConfig.PoolSize = utils.GetEnvInt("REDIS_POOL_SIZE", 10)
	LightRedisConfig.MinIdleConns = utils.GetEnvInt("REDIS_MIN_IDLE_CONNS", 5)

	if LightRedisConfig.Enable {
		initRedis()
	}
}
