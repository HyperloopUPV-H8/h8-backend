package internals

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/models"
)

func encodeNext(writer io.Writer, value any) {
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		log.Fatalf("encoder: encodeNext: %s\n", err)
	}
}

func EncodeID(writer io.Writer, id uint16) {
	encodeNext(writer, id)
}

func EncodeEnum(writer io.Writer, enum models.Enum, value string) {
	encodeNext(writer, enum.GetNumericValue(value))
}

func EncodeBool(writer io.Writer, value bool) {
	encodeNext(writer, value)
}

func EncodeString(writer io.Writer, value string) {
	encodeNext(writer, []byte(value))
}

func EncodeNumber(writer io.Writer, kind string, value float64) {
	switch kind {
	case "uint8":
		encodeNext(writer, uint8(value))
	case "uint16":
		encodeNext(writer, uint16(value))
	case "uint32":
		encodeNext(writer, uint32(value))
	case "uint64":
		encodeNext(writer, uint64(value))
	case "int8":
		encodeNext(writer, int8(value))
	case "int16":
		encodeNext(writer, int16(value))
	case "int32":
		encodeNext(writer, int32(value))
	case "int64":
		encodeNext(writer, int64(value))
	case "float32":
		encodeNext(writer, float32(value))
	default:
		encodeNext(writer, float64(value))
	}
}
