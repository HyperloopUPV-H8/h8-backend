package serde

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type ID = uint16
type Numeric = float64
type ValueType = string
type EnumVariant = string
type Enum = map[uint8]EnumVariant

func decodeRaw[T any](bytes io.Reader) (value T) {
	err := binary.Read(bytes, binary.LittleEndian, &value)
	if err != nil {
		panic(err)
	}

	return
}

func DecodeNumber(numType ValueType, bytes io.Reader) Numeric {
	switch numType {
	case "uint8":
		return Numeric(decodeRaw[uint8](bytes))
	case "uint16":
		return Numeric(decodeRaw[uint16](bytes))
	case "uint32":
		return Numeric(decodeRaw[uint32](bytes))
	case "uint64":
		return Numeric(decodeRaw[uint64](bytes))
	case "int8":
		return Numeric(decodeRaw[int8](bytes))
	case "int16":
		return Numeric(decodeRaw[int16](bytes))
	case "int32":
		return Numeric(decodeRaw[int32](bytes))
	case "int64":
		return Numeric(decodeRaw[int64](bytes))
	case "float32":
		return Numeric(decodeRaw[float32](bytes))
	case "float64":
		return Numeric(decodeRaw[float64](bytes))
	default:
		panic(fmt.Sprintf("Expected numeric type, got %s", numType))
	}
}

func DecodeBool(bytes io.Reader) bool {
	return decodeRaw[bool](bytes)
}

func DecodeEnum(enum Enum, bytes io.Reader) EnumVariant {
	value, exists := enum[decodeRaw[uint8](bytes)]
	if !exists {
		log.Fatalln("decode enum: expected value, got nil")
	}
	return value
}

func DecodeString(bytes io.Reader) string {
	line, err := bufio.NewReader(bytes).ReadString('\n')
	if err != nil {
		log.Fatalln("decode string:", err)
	}
	return line
}

func DecodeID(bytes io.Reader) ID {
	return ID(decodeRaw[uint16](bytes))
}
