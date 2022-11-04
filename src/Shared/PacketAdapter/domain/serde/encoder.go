package serde

// Currently this file is not used because currently PacketParser doesn't need to encode orders

import (
	"encoding/binary"
	"fmt"
	"io"
)

func encodeRaw(value any, bytes io.Writer) {
	err := binary.Write(bytes, binary.BigEndian, value)
	if err != nil {
		panic(err)
	}
}

func EncodeNumber(numType ValueType, value Numeric, bytes io.Writer) {
	switch numType {
	case "uint8":
		encodeRaw(uint8(value), bytes)
	case "uint16":
		encodeRaw(uint16(value), bytes)
	case "uint32":
		encodeRaw(uint32(value), bytes)
	case "uint64":
		encodeRaw(uint64(value), bytes)
	case "int8":
		encodeRaw(int8(value), bytes)
	case "int16":
		encodeRaw(int16(value), bytes)
	case "int32":
		encodeRaw(int32(value), bytes)
	case "int64":
		encodeRaw(int64(value), bytes)
	case "float32":
		encodeRaw(float32(value), bytes)
	case "float64":
		encodeRaw(value, bytes)
	default:
		panic(fmt.Sprintf("Expected numeric type, got %s", numType))
	}
}

func EncodeBool(value bool, bytes io.Writer) {
	encodeRaw(value, bytes)
}

func EncodeEnum(enum Enum, value EnumVariant, bytes io.Writer) {
	for code, v := range enum {
		if v == value {
			encodeRaw(code, bytes)
			return
		}
	}
}

func EncodeID(id uint16, bytes io.Writer) {
	encodeRaw(id, bytes)
}
