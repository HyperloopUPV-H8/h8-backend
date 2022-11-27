package decoder

import (
	"encoding/binary"
	"io"
	"log"
)

func decodeNext[T any](reader io.Reader) (value T) {
	if err := binary.Read(reader, binary.LittleEndian, value); err != nil {
		log.Fatalf("packet parser: decode next: %s\n", err)
	}
	return value
}
