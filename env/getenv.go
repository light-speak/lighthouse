package env

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	loadedEnv = false
)

func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func Init(path *string) (err error) {
	if loadedEnv {
		return nil
	}
	p := ""
	if path == nil {
		p = "../.env"
	} else {
		p = *path
	}
	loadedEnv = true
	return godotenv.Load(p)
}
