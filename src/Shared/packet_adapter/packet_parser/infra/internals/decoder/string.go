package decoder

import (
	"bufio"
	"io"
	"log"
)

func String(reader io.Reader) string {
	line, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		log.Fatalln("decode string:", err)
	}
	return line
}
