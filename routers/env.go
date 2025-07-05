package routers

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
)

// JWT_SECRET is the secret key for the JWT token
type middlewareConfig struct {
	// JWT SECRET
	JWT_SECRET string
	// HeartbeatPath is the path for the heartbeat endpoint
	HeartbeatPath string
	// CompressLevel is the level of compression for the response
	CompressLevel int
	// Timeout is the timeout for the request
	Timeout time.Duration
	// Throttle is the throttle for the request
	Throttle int
}

var Config *middlewareConfig

func init() {
	Config = &middlewareConfig{
		JWT_SECRET:    "IWY@*3JUI#d309HhefzX2WpLtPKtD!hn",
		HeartbeatPath: "/health",
		CompressLevel: 5,
		Timeout:       10 * time.Second,
		Throttle:      100,
	}

	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	Config.JWT_SECRET = utils.GetEnv("JWT_SECRET", Config.JWT_SECRET)
	Config.HeartbeatPath = utils.GetEnv("MID_HEARTBEAT_PATH", Config.HeartbeatPath)
	Config.CompressLevel = utils.GetEnvInt("MID_COMPRESS_LEVEL", Config.CompressLevel)
	Config.Timeout = time.Duration(utils.GetEnvInt("MID_TIMEOUT", 30)) * time.Second
	Config.Throttle = utils.GetEnvInt("MID_THROTTLE", Config.Throttle)
}
