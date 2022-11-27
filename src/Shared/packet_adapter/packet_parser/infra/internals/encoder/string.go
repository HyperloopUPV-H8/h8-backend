package encoder

import (
	"io"
	"log"
)

func String(writer io.Writer, value string) {
	if _, err := writer.Write([]byte(value)); err != nil {
		log.Fatalf("packet parser: encode string: %s\n", err)
	}
}
