package decoder

import "io"

func Bool(reader io.Reader) bool {
	return decodeNext[bool](reader)
}
