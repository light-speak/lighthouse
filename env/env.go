package env

import (
	"os"
	"strconv"
)

// GetEnv Get the value of the environment variable
func GetEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetEnvInt64 Get the value of the environment variable as an int64
func GetEnvInt64(key string, defaultValue ...int64) int64 {
	value, err := strconv.ParseInt(GetEnv(key), 10, 64)
	if err != nil {
		value = 0
	}
	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetEnvInt Get the value of the environment variable as an int
func GetEnvInt(key string, defaultValue ...int) int {
	value, err := strconv.Atoi(GetEnv(key))
	if err != nil {
		value = 0
	}
	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetEnvBool Get the value of the environment variable as a bool
func GetEnvBool(key string, defaultValue ...bool) bool {
	value, err := strconv.ParseBool(GetEnv(key))
	if err != nil {
		value = false
	}
	return value
}

// GetEnvFloat64 Get the value of the environment variable as a float64
func GetEnvFloat64(key string, defaultValue ...float64) float64 {
	value, err := strconv.ParseFloat(GetEnv(key), 64)
	if err != nil {
		value = 0
	}
	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
