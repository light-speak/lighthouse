package env

import "os"

// Getenv 获取Env Value，带默认值
func Getenv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}
