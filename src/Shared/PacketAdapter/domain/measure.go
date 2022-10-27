package domain

import (
	"io"

	value "github.com/HyperloopUPV-H8/Backend-H8/..."
)

type Measure struct {
	name      Name
	valueType ValueType
}

func (measure Measure) Decode(enums map[Name]Enum, reader io.Reader) value.Value {
	if enum, exists := enums[Name(measure.valueType)]; exists {
		return DecodeEnum(enum, reader)
	} else if measure.valueType == "bool" {
		return DecodeBool(reader)
	} else {
		return DecodeNumber(measure.valueType, reader)
	}
}
