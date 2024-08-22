package parser

import "strings"

func ParseStringArray(arrayStr string) ([]string, error) {
	trimmed := strings.Trim(arrayStr, "{}")

	items := strings.Split(trimmed, ",")

	return items, nil
}
