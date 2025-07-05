package utils

import (
	"os"
	"strconv"
	"strings"
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
	value := GetEnv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return result
}

// GetEnvInt Get the value of the environment variable as an int
func GetEnvInt(key string, defaultValue ...int) int {
	value := GetEnv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return result
}

// GetEnvBool Get the value of the environment variable as a bool
func GetEnvBool(key string, defaultValue ...bool) bool {
	value := GetEnv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return result
}

// GetEnvFloat64 Get the value of the environment variable as a float64
func GetEnvFloat64(key string, defaultValue ...float64) float64 {
	value := GetEnv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return result
}

// GetEnvArray Get the value of the environment variable as a string array
func GetEnvArray(key string, sep string, defaultValue ...[]string) []string {
	value := GetEnv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []string{}
	}
	return strings.Split(value, sep)
}

// GetEnvInt64Array Get the value of the environment variable as an int64 array
func GetEnvInt64Array(key string, sep string, defaultValue ...[]int64) []int64 {
	strArr := GetEnvArray(key, sep)
	if len(strArr) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []int64{}
	}
	var result []int64
	for _, str := range strArr {
		if value, err := strconv.ParseInt(str, 10, 64); err == nil {
			result = append(result, value)
		}
	}
	return result
}

// GetEnvIntArray Get the value of the environment variable as an int array
func GetEnvIntArray(key string, sep string, defaultValue ...[]int) []int {
	strArr := GetEnvArray(key, sep)
	if len(strArr) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []int{}
	}
	var result []int
	for _, str := range strArr {
		if value, err := strconv.Atoi(str); err == nil {
			result = append(result, value)
		}
	}
	return result
}

// GetEnvFloat64Array Get the value of the environment variable as a float64 array
func GetEnvFloat64Array(key string, sep string, defaultValue ...[]float64) []float64 {
	strArr := GetEnvArray(key, sep)
	if len(strArr) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []float64{}
	}
	var result []float64
	for _, str := range strArr {
		if value, err := strconv.ParseFloat(str, 64); err == nil {
			result = append(result, value)
		}
	}
	return result
}
