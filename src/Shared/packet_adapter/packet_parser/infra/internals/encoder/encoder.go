package encoder

import (
	"encoding/binary"
	"io"
	"log"
)

func encodeNext(writer io.Writer, value any) {
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		log.Fatalf("packet parser: encode next: %s\n", err)
	}
}
