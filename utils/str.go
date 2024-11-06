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

// CamelCase 驼峰命名
func CamelCase(str string) string {
	// 如果已经是驼峰命名，直接返回
	if !strings.Contains(str, "_") {
		// 确保首字母小写
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
