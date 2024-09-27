package utils

import "strings"

func LcFirst(str string) string {
	if str == "" {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

func ToLower(str string) string {
	return strings.ToLower(str)
}

func Pluralize(word string) string {
	if word == "" {
		return word
	}

	// 简单的复数规则
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return strings.ToLower(word + "es")
	}

	if strings.HasSuffix(word, "y") {
		if len(word) > 1 && !isVowel(rune(word[len(word)-2])) {
			return strings.ToLower(word[:len(word)-1] + "ies")
		}
		return strings.ToLower(word + "s")
	}

	return strings.ToLower(word + "s")
}

func isVowel(c rune) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}

func UcFirst(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func UcFirstWithID(str string) string {
	if str == "" {
		return str
	}
	str = strings.ToUpper(str[:1]) + str[1:]
	return strings.ReplaceAll(str, "Id", "ID")
}
