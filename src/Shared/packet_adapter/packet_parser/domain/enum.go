package domain

import (
	"regexp"
	"strings"
)

type enumVariant = string
type Enum = map[uint8]enumVariant

var enumExp = regexp.MustCompile(`(?i)^ENUM\((\w+(?:,\w+)*)\)$`)
var itemsExp = regexp.MustCompile(`(\w+),?`)

func NewEnum(enumString string) Enum {
	matches := getEnumMatches(enumString)
	return parseEnum(matches)
}

func parseEnum(matches [][]string) Enum {
	variants := make(map[uint8]enumVariant, len(matches))
	for i, match := range matches {
		variants[uint8(i)] = enumVariant(match[1])
	}
	return Enum(variants)
}

func getEnumMatches(enumString string) [][]string {
	itemList := enumExp.FindStringSubmatch(removeWhitespace(enumString))[1]
	return itemsExp.FindAllStringSubmatch(itemList, -1)
}

func removeWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}

func IsEnum(valueType string) bool {
	return enumExp.MatchString(removeWhitespace(valueType))
}
