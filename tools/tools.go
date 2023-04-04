package tools

import "strings"

func ReplaceQuotes(str string) string {
	return strings.ReplaceAll(str, "'", "\"")
}
