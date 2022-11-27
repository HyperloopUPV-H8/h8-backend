package encoder

import "io"

func Bool(writer io.Writer, value bool) {
	encodeNext(writer, value)
}
