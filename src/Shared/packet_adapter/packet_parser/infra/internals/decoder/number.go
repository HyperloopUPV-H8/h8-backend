package decoder

import (
	"io"
)

func ID(reader io.Reader) uint16 {
	return decodeNext[uint16](reader)
}

func Number(reader io.Reader, kind string) float64 {
	switch kind {
	case "uint8":
		return float64(decodeNext[uint8](reader))
	case "uint16":
		return float64(decodeNext[uint16](reader))
	case "uint32":
		return float64(decodeNext[uint32](reader))
	case "uint64":
		return float64(decodeNext[uint64](reader))
	case "int8":
		return float64(decodeNext[int8](reader))
	case "int16":
		return float64(decodeNext[int16](reader))
	case "int32":
		return float64(decodeNext[int32](reader))
	case "int64":
		return float64(decodeNext[int64](reader))
	case "float32":
		return float64(decodeNext[float32](reader))
	case "float64":
		return decodeNext[float64](reader)
	}

	panic("packet parser: decode number: invalid kind " + kind + "\n")
}
