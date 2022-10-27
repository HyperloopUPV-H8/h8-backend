package value

import "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value/number"

type Value interface {
	ToDisplayString() string
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

func isNumber(valueType string) bool {
	switch valueType {
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64":
		return true
	default:
		return false
	}
}
// TODO: Esta función esta hecha pero no se donde colocarla, por ahora dejo este placeholder
func isEnum(va string) bool {
	return true
}
