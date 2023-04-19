package message

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
)

type Parser struct {
	config Config
}

func NewParser(config Config) Parser {
	return Parser{config: config}
}

func (parser Parser) Decode(id uint16, data []byte) (packet.Payload, error) {
	switch id {
	case parser.config.FaultId:
		protection, err := NewProtection("fault", data)
		return Payload{Data: protection, raw: data}, err
	case parser.config.WarningId:
		protection, err := NewProtection("warning", data)
		return Payload{Data: protection, raw: data}, err
	case parser.config.BlcuAckId:
		return Payload{Data: BlcuAck{raw: data}, raw: data}, nil
	default:
		return nil, fmt.Errorf("unknown message id %d", id)
	}
}
