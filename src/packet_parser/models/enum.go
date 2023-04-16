package models

import (
	"fmt"
	"strings"
)

type Enum map[uint8]string

func (enum Enum) GetNumericValue(value string) uint8 {
	for repr, variant := range enum {
		if value == variant {
			return repr
		}
	}
	panic(fmt.Sprintf("Enum: GetRepr: failed to get representation for %s", value))
}

func GetEnum(literal string) Enum {
	enum := make(Enum)
	for i, variant := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(literal, "ENUM("), ")"), ",") {
		enum[uint8(i)] = strings.Trim(variant, " ")
	}
	return enum
}
