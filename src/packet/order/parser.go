package order

import (
	"bytes"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/parsers"
)

type Parser struct {
	value    *parsers.ValueParser
	bitarray *parsers.BitarrayParser
}

func NewParser(valueParser *parsers.ValueParser, bitarrayParser *parsers.BitarrayParser) Parser {
	return Parser{value: valueParser, bitarray: bitarrayParser}
}

func (parser Parser) Decode(id uint16, data []byte) (packet.Payload, error) {
	reader := bytes.NewReader(data)
	values, err := parser.value.Decode(id, reader)
	if err != nil {
		return Payload{}, err
	}

	enabled, err := parser.bitarray.Decode(id, reader)
	if err != nil {
		return Payload{}, err
	}

	return Payload{Values: values, Enabled: enabled, raw: data}, nil
}

func (parser Parser) Encode(id uint16, payload packet.Payload) ([]byte, error) {
	orderPayload, ok := payload.(Payload)
	if !ok {
		return nil, fmt.Errorf("invalid order payload type %T", payload)
	}

	buf := bytes.NewBuffer(nil)

	err := parser.value.Encode(id, orderPayload.Values, buf)
	if err != nil {
		return nil, err
	}

	err = parser.bitarray.Encode(id, orderPayload.Enabled, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
