package utils

import (
	"strings"
	"unicode"
)

var specialWords = map[string]string{
	"id":    "ID",
	"ip":    "IP",
	"url":   "URL",
	"uri":   "URI",
	"api":   "API",
	"uuid":  "UUID",
	"html":  "HTML",
	"xml":   "XML",
	"json":  "JSON",
	"yaml":  "YAML",
	"css":   "CSS",
	"sql":   "SQL",
	"http":  "HTTP",
	"https": "HTTPS",
	"ftp":   "FTP",
	"ssh":   "SSH",
	"ssl":   "SSL",
	"tcp":   "TCP",
	"udp":   "UDP",
	"gui":   "GUI",
	"ui":    "UI",
	"cdn":   "CDN",
	"dns":   "DNS",
	"cpu":   "CPU",
	"gpu":   "GPU",
	"ram":   "RAM",
	"sdk":   "SDK",
	"jwt":   "JWT",
	"oauth": "OAuth",
}

// LcFirst 首字母小写
func LcFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

// 转小写
func Lc(str string) string {
	return strings.ToLower(str)
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
	var lastIsUpper bool
	var currentWord []rune

	for i, char := range str {
		isUpper := unicode.IsUpper(char)

		if isUpper {
			if !lastIsUpper && i > 0 {
				result = append(result, string(currentWord))
				currentWord = []rune{}
			} else if lastIsUpper && i+1 < len(str) {
				nextChar := rune(str[i+1])
				if !unicode.IsUpper(nextChar) && !unicode.IsNumber(nextChar) {
					result = append(result, string(currentWord))
					currentWord = []rune{}
				}
			}
		}

		lastIsUpper = isUpper
		currentWord = append(currentWord, unicode.ToLower(char))
	}

	if len(currentWord) > 0 {
		result = append(result, string(currentWord))
	}

	return strings.Join(result, "_")
}

// CamelCase 驼峰命名（增强版，支持特殊单词）
func CamelCase(str string) string {
	// 如果不包含下划线，且不是特殊单词，直接返回首字母小写
	if !strings.Contains(str, "_") {
		return LcFirst(str)
	}

	str = strings.TrimLeft(str, "_")
	str = strings.TrimRight(str, "_")

	parts := strings.Split(str, "_")
	var result []string

	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result = append(result, strings.ToLower(part))
		} else {
			result = append(result, UcFirst(strings.ToLower(part)))
		}
	}

	return strings.Join(result, "")
}

// CamelCaseWithSpecial 驼峰命名（特殊单词会保持大写）
func CamelCaseWithSpecial(str string) string {
	// Check if it's a special word first
	if v, ok := specialWords[strings.ToLower(str)]; ok {
		return v
	}

	// If no underscore and not a special word
	if !strings.Contains(str, "_") {
		// Split the string into parts based on uppercase letters
		var parts []string
		var currentPart string

		for i, char := range str {
			if i > 0 && unicode.IsUpper(char) {
				if len(currentPart) > 0 {
					parts = append(parts, currentPart)
				}
				currentPart = string(char)
			} else {
				currentPart += string(char)
			}
		}
		if len(currentPart) > 0 {
			parts = append(parts, currentPart)
		}

		// Process each part
		var result []string
		for _, part := range parts {
			if v, ok := specialWords[strings.ToLower(part)]; ok {
				result = append(result, v)
			} else if strings.ToLower(part) == "id" {
				result = append(result, "ID")
			} else {
				result = append(result, UcFirst(strings.ToLower(part)))
			}
		}
		return strings.Join(result, "")
	}

	str = strings.TrimLeft(str, "_")
	str = strings.TrimRight(str, "_")

	parts := strings.Split(str, "_")
	var result []string

	for _, part := range parts {
		if part == "" {
			continue
		}
		// Check if part is a special word
		if v, ok := specialWords[strings.ToLower(part)]; ok {
			result = append(result, v)
		} else if strings.ToLower(part) == "id" {
			result = append(result, "ID")
		} else {
			// Capitalize first letter of each part
			result = append(result, UcFirst(strings.ToLower(part)))
		}
	}

	return strings.Join(result, "")
}

func StrPtr(str string) *string {
	return &str
}

// Irregular plural forms
var irregulars = map[string]string{
	"man":    "men",
	"woman":  "women",
	"child":  "children",
	"tooth":  "teeth",
	"foot":   "feet",
	"mouse":  "mice",
	"person": "people",
}

// Pluralize function to convert a word to its plural form
func Pluralize(word string) string {
	// Check for irregular forms
	if plural, exists := irregulars[word]; exists {
		return plural
	}

	// General rules
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	} else if len(word) > 1 && strings.HasSuffix(word, "y") && !isVowel(rune(word[len(word)-2])) {
		return word[:len(word)-1] + "ies"
	} else if strings.HasSuffix(word, "y") {
		return word + "s"
	}

	return word + "s"
}

// Helper function to check if a character is a vowel
func isVowel(c rune) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}

func IsInternalType(name string) bool {
	return len(name) >= 2 && name[:2] == "__"
}

func Able(name string) string {
	// Handle special cases ending in 'e'
	if strings.HasSuffix(name, "e") {
		return name[:len(name)-1] + "able"
	}

	// Handle special cases ending in 'y'
	if strings.HasSuffix(name, "y") {
		return name[:len(name)-1] + "iable"
	}

	return name + "able"
}

// CamelColon 将驼峰命名转换为路径格式 (例如: ProcessTask -> process::task)
func CamelColon(str string) string {
	// 先转换为蛇形命名
	snake := SnakeCase(str)
	// 将下划线替换为双冒号
	return strings.ReplaceAll(snake, "_", "::")
}

// PascalCase 帕斯卡命名（首字母大写的驼峰）
func PascalCase(str string) string {
	return UcFirst(CamelCase(str))
}
