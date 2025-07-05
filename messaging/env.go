package messaging

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
)

var cfg Config

type Driver string

const (
	DriverNats  Driver = "nats"
	DriverKafka Driver = "kafka"
)

type Config struct {
	Driver     Driver
	URL        string
	InstanceID string
}

func init() {
	cfg = Config{
		Driver:     DriverNats,
		URL:        "localhost:4222",
		InstanceID: "hostname",
	}

	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	cfg.Driver = Driver(utils.GetEnv("MESSAGING_DRIVER", string(cfg.Driver)))
	cfg.URL = utils.GetEnv("MESSAGING_URL", cfg.URL)
	cfg.InstanceID = utils.GetEnv("HOSTNAME", cfg.InstanceID)

	if cfg.Driver == "" {
		log.Fatal("MESSAGING_DRIVER is not set")
	}
}
