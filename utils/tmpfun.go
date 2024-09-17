package utils

import "strings"

func LcFirst(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToLower(string(str[0])) + str[1:]
}

func UcFirst(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func UcFirstWithID(str string) string {
	if str == "" {
		return ""
	}
	str = strings.ToUpper(string(str[0])) + str[1:]

	// 替换特定部分为大写
	str = strings.ReplaceAll(str, "Id", "ID")

	return str
}
