package utils

import "strings"

func LcFirst(str string) string {
	if str == "" {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
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
