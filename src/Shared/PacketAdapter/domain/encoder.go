package domain

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

func EncodeNumber(numType ValueType, value float64, bytes io.Writer) {
	switch numType {
	case "Uint8":
		encodeRaw(uint8(value), bytes)
	case "Uint16":
		encodeRaw(uint16(value), bytes)
	case "Uint32":
		encodeRaw(uint32(value), bytes)
	case "Uint64":
		encodeRaw(uint64(value), bytes)
	case "Int8":
		encodeRaw(int8(value), bytes)
	case "Int16":
		encodeRaw(int16(value), bytes)
	case "Int32":
		encodeRaw(int32(value), bytes)
	case "Int64":
		encodeRaw(int64(value), bytes)
	case "Float32":
		encodeRaw(float32(value), bytes)
	case "Float64":
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
