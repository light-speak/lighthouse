package queue

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
)

type QueueConfig struct {
	Enable   bool
	Host     string
	Port     string
	Password string
	DB       int
}

var LightQueueConfig *QueueConfig

func init() {
	LightQueueConfig = &QueueConfig{
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

	LightQueueConfig.Enable = utils.GetEnvBool("QUEUE_ENABLE", LightQueueConfig.Enable)
	if LightQueueConfig.Enable {
		LightQueueConfig.Host = utils.GetEnv("QUEUE_REDIS_HOST", LightQueueConfig.Host)
		LightQueueConfig.Port = utils.GetEnv("QUEUE_REDIS_PORT", LightQueueConfig.Port)
		LightQueueConfig.Password = utils.GetEnv("QUEUE_REDIS_PASSWORD", LightQueueConfig.Password)
		LightQueueConfig.DB = utils.GetEnvInt("QUEUE_REDIS_DB", LightQueueConfig.DB)
	}

	if LightQueueConfig.Enable && (LightQueueConfig.Host == "" || LightQueueConfig.Port == "" || LightQueueConfig.Password == "" || LightQueueConfig.DB == 0) {
		log.Println("Queue config is invalid, please check the .env file")
		panic("Queue config is invalid, please check the .env file")
	}
}
