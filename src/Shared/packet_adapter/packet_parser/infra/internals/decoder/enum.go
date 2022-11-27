package decoder

import (
	"io"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

func Enum(reader io.Reader, enum domain.Enum) string {
	k := decodeNext[uint8](reader)
	variant, exists := enum[k]
	if !exists {
		log.Fatalf("packet parser: decode enum: got invalid variant %d (%d)", k, len(enum)-1)
	}
	return variant
}
