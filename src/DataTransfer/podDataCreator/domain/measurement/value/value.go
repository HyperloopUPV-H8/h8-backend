package value

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value/number"
)

type Value interface {
	//current() Value
	//fromBuffer(b *bytes.Buffer)
	//toString() string
	ToDisplayString() string
}

func NewDefault(valueType string) Value {
	switch valueType {
	case "bool":
		return new(Bool)
	default:
		if isEnum(valueType) {
			return new(String)
		} else if isNumber(valueType) {
			return new(number.Number)
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

func isEnum(va string) bool {
	return true
}

// func New(valueType string) Value {
// 	switch valueType {
// 	case "uint8":
// 		return new(UInt8)
// 	case "uint16":
// 		return new(UInt16)
// 	case "uint32":
// 		return new(UInt32)
// 	case "uint64":
// 		return new(UInt64)
// 	case "int8":
// 		return new(Int8)
// 	case "int16":
// 		return new(Int16)
// 	case "int32":
// 		return new(Int32)
// 	case "int64":
// 		return new(Int64)
// 	case "float32":
// 		return new(Float32)
// 	case "float64":
// 		return new(Float64)
// 	case "bool":
// 		return new(Bool)
// 	default:
// 		if isEnum(valueType) {
// 			return newEnum(valueType)
// 		} else {
// 			panic("Invalid type")
// 		}
// 	}
// }
