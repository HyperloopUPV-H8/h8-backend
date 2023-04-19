package data

import (
	"bytes"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/parsers"
)

type Parser struct {
	value *parsers.ValueParser
}

func NewParser(valueParser *parsers.ValueParser) Parser {
	return Parser{value: valueParser}
}

func (parser Parser) Decode(id uint16, data []byte) (packet.Payload, error) {
	values, err := parser.value.Decode(id, bytes.NewReader(data))
	return Payload{Values: values}, err
}
