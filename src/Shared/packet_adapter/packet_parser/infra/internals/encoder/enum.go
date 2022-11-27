package encoder

import (
	"io"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

func Enum(writer io.Writer, enum domain.Enum, value string) {
	val, exists := enum.Find(value)
	if !exists {
		log.Fatalf("packet parser: encode enum: invalid value %s\n", value)
	}
	encodeNext(writer, val)
}
