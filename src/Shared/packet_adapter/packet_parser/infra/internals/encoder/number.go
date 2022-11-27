package encoder

import (
	"io"
	"log"
)

func ID(writer io.Writer, id uint16) {
	encodeNext(writer, id)
}

func Number(writer io.Writer, value float64, kind string) {
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
	case "float64":
		encodeNext(writer, value)
	default:
		log.Fatalf("packet parser: encode number: invalid kind %s\n", kind)
	}
}
