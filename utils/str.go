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

// Irregular plural forms
var irregulars = map[string]string{
	"man":   "men",
	"woman": "women",
	"child": "children",
	"tooth": "teeth",
	"foot":  "feet",
	"mouse": "mice",
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