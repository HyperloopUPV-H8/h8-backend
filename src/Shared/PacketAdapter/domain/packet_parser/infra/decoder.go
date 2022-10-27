package infra

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
)

func decodeRaw[T any](bytes io.Reader) (value T) {
	err := binary.Read(bytes, binary.BigEndian, &value)
	if err != nil {
		panic(err)
	}

	return
}

func DecodeNumber(numType domain.ValueType, bytes io.Reader) domain.Numeric {
	switch numType {
	case "uint8":
		return domain.Numeric(decodeRaw[uint8](bytes))
	case "uint16":
		return domain.Numeric(decodeRaw[uint16](bytes))
	case "uint32":
		return domain.Numeric(decodeRaw[uint32](bytes))
	case "uint64":
		return domain.Numeric(decodeRaw[uint64](bytes))
	case "int8":
		return domain.Numeric(decodeRaw[int8](bytes))
	case "int16":
		return domain.Numeric(decodeRaw[int16](bytes))
	case "int32":
		return domain.Numeric(decodeRaw[int32](bytes))
	case "int64":
		return domain.Numeric(decodeRaw[int64](bytes))
	case "float32":
		return domain.Numeric(decodeRaw[float32](bytes))
	case "float64":
		return domain.Numeric(decodeRaw[float64](bytes))
	default:
		panic(fmt.Sprintf("Expected numeric type, got %s", numType))
	}
}

func DecodeBool(bytes io.Reader) bool {
	return decodeRaw[bool](bytes)
}

func DecodeEnum(enum domain.Enum, bytes io.Reader) domain.EnumVariant {
	return enum[decodeRaw[uint8](bytes)]
}

func DecodeID(bytes io.Reader) domain.ID {
	return domain.ID(decodeRaw[uint16](bytes))
}
