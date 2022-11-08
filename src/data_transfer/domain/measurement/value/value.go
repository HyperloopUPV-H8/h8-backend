package value

import (
	"regexp"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain/measurement/value/number"
)

type Value interface {
	ToPodUnitsString() string
	ToDisplayUnitsString() string
	GetPodUnits() string
	GetDisplayUnits() string
	Update(newValue any)
}

func NewDefault(valueType string, podUnits string, displayUnits string) Value {
	switch valueType {
	case "bool":
		return new(Bool)
	default:
		if isEnum(valueType) {
			return new(String)
		} else if isNumber(valueType) {
			return number.NewNumber(podUnits, displayUnits)
		} else {
			panic("Invalid type")
		}
	}
}

var numExp = regexp.MustCompile(`(?i)U?INT(?:8|16|32|64)|float(?:32|64)`)

func isNumber(valueType string) bool {
	return numExp.MatchString(valueType)
}

var enumExp = regexp.MustCompile(`(?i)^ENUM\((\w+(?:,\w+)*)\)$`)

func isEnum(valueType string) bool {
	return enumExp.MatchString(removeWhitespace(valueType))
}

func removeWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}
