package internals

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/models"
	trace "github.com/rs/zerolog/log"
)

func decodeNext[T any](reader io.Reader) (value T) {
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		trace.Fatal().Stack().Err(err).Msg("")
		return
	}
	return value
}

func DecodeID(reader io.Reader) uint16 {
	return decodeNext[uint16](reader)
}

func DecodeEnum(reader io.Reader, enum models.Enum) string {
	return enum[decodeNext[uint8](reader)]
}

func DecodeBool(reader io.Reader) bool {
	return decodeNext[bool](reader)
}

func DecodeString(reader io.Reader) string {
	str, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("")
		return ""
	}
	return str
}

func DecodeNumber(reader io.Reader, kind string) float64 {
	switch kind {
	case "uint8":
		return float64(decodeNext[uint8](reader))
	case "uint16":
		return float64(decodeNext[uint16](reader))
	case "uint32":
		return float64(decodeNext[uint32](reader))
	case "uint64":
		return float64(decodeNext[uint64](reader))
	case "int8":
		return float64(decodeNext[int8](reader))
	case "int16":
		return float64(decodeNext[int16](reader))
	case "int32":
		return float64(decodeNext[int32](reader))
	case "int64":
		return float64(decodeNext[int64](reader))
	case "float32":
		return float64(decodeNext[float32](reader))
	default:
		return decodeNext[float64](reader)
	}
}
