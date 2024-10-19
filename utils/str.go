package utils

import (
	"strings"
	"unicode"
)

// LcFirst 首字母小写
func LcFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

// UcFirst 首字母大写
func UcFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// SnakeCase 下划线命名
func SnakeCase(str string) string {
	var result []string
	for i, char := range str {
		if unicode.IsUpper(char) && i > 0 {
			result = append(result, "_")
		}
		result = append(result, string(unicode.ToLower(char)))
	}
	return strings.Join(result, "")
}

// CamelCase 驼峰命名
func CamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i := range parts {
		if i > 0 {
			parts[i] = UcFirst(parts[i])
		}
	}
	return strings.Join(parts, "")
}

func StrPtr(str string) *string {
	return &str
}
