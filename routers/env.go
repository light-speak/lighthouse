package routers

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/light-speak/lighthouse/utils"
)

// JWT_SECRET is the secret key for the JWT token
type middlewareConfig struct {
	// JWT SECRET
	JWT_SECRET string
	// HeartbeatPath is the path for the liveness endpoint
	HeartbeatPath string
	// ReadinessPath is the path for the readiness endpoint
	ReadinessPath string
	// CompressLevel is the level of compression for the response
	CompressLevel int
	// Timeout is the timeout for the request
	Timeout time.Duration
	// Throttle is the throttle for the request
	Throttle int

	CORSAllowOrigins []string
	CORSAllowMethods []string
	CORSAllowHeaders []string
}

var Config *middlewareConfig

func init() {
	Config = &middlewareConfig{
		JWT_SECRET:       "IWY@*3JUI#d309HhefzX2WpLtPKtD!hn",
		HeartbeatPath:    "/health",
		ReadinessPath:    "/ready",
		CompressLevel:    5,
		Timeout:          10 * time.Second,
		Throttle:         100,
		CORSAllowOrigins: []string{"*"},
		CORSAllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"},
		CORSAllowHeaders: []string{"*"},
	}

	if curPath, err := os.Getwd(); err == nil {
		err = godotenv.Load(filepath.Join(curPath, ".env"))
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	Config.JWT_SECRET = utils.GetEnv("JWT_SECRET", Config.JWT_SECRET)
	Config.HeartbeatPath = utils.GetEnv("MID_HEARTBEAT_PATH", Config.HeartbeatPath)
	Config.ReadinessPath = utils.GetEnv("MID_READINESS_PATH", Config.ReadinessPath)
	Config.CompressLevel = utils.GetEnvInt("MID_COMPRESS_LEVEL", Config.CompressLevel)
	Config.Timeout = time.Duration(utils.GetEnvInt("MID_TIMEOUT", 30)) * time.Second
	Config.Throttle = utils.GetEnvInt("MID_THROTTLE", Config.Throttle)
	if origins := utils.GetEnv("CORS_ALLOW_ORIGINS", ""); origins != "" {
		Config.CORSAllowOrigins = strings.Split(origins, ",")
	}
	if methods := utils.GetEnv("CORS_ALLOW_METHODS", ""); methods != "" {
		Config.CORSAllowMethods = strings.Split(methods, ",")
	}
	if headers := utils.GetEnv("CORS_ALLOW_HEADERS", ""); headers != "" {
		Config.CORSAllowHeaders = strings.Split(headers, ",")
	}
}
