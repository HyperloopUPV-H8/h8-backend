package value

import (
	"bytes"
	"encoding/binary"
)

type Value interface {
	current() Value
	fromBuffer(b bytes.Buffer)
	//toBytes() []byte
}

func New(valueType string) Value {
	switch valueType {
	case "uint8":
		return new(UInt8)
	case "uint16":
		return new(UInt16)
	case "uint32":
		return new(UInt32)
	case "uint64":
		return new(UInt64)
	case "int8":
		return new(Int8)
	case "int16":
		return new(Int16)
	case "int32":
		return new(Int32)
	case "int64":
		return new(Int64)
	case "float32":
		return new(Float32)
	case "float64":
		return new(Float64)
	case "bool":
		return new(Bool)
	default:
		panic("Invalid type")
	}
}

func numberFromBuffer(b bytes.Buffer, n int) uint64 {
	bytes := b.Next(n)
	return binary.BigEndian.Uint64(bytes)
}
