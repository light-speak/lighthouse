package env

import "os"
import "github.com/joho/godotenv"

func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func Init(path *string) (err error) {
	p := ""
	if path == nil {
		p = "../.env"
	} else {
		p = *path
	}
	return godotenv.Load(p)
}
