package domain

import (
	"strings"
)

type Enum map[uint8]string

func (enum Enum) Find(value string) (uint8, bool) {
	for val, repr := range enum {
		if repr == value {
			return val, true
		}
	}
	return 0, false
}

func NewEnum(literal string) Enum {
	enum := make(map[uint8]string)
	for i, variant := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(removeWhitespace(literal), "ENUM("), ")"), ",") {
		enum[uint8(i)] = variant
	}
	return enum
}

func removeWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}

func IsEnum(literal string) bool {
	return strings.HasPrefix(literal, "ENUM")
}
