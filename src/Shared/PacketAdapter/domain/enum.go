package domain

import (
	"regexp"
	"strings"
)

type Enum map[uint8]EnumVariant

var enumExpr = regexp.MustCompile(`^ENUM\((\w+(?:,\w+)*)\)$`)
var itemsExpr = regexp.MustCompile(`(\w+),?`)

func NewEnum(enumString string) Enum {
	matches := getEnumMatches(enumString)
	return parseEnum(matches)
}

func parseEnum(matches [][]string) Enum {
	variants := make(map[uint8]EnumVariant, len(matches))
	for i, match := range matches {
		variants[uint8(i)] = EnumVariant(match[1])
	}
	return Enum(variants)
}

func getEnumMatches(enumString string) [][]string {
	itemList := enumExpr.FindStringSubmatch(removeWhitespace(enumString))[1]
	return itemsExpr.FindAllStringSubmatch(itemList, -1)
}

func removeWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}
