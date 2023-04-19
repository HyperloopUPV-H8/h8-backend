package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Parser struct {
	idToKind map[uint16]Kind
	decoders map[Kind]Decoder
	encoders map[Kind]Encoder
}

type Decoder interface {
	Decode(id uint16, data []byte) (Payload, error)
}

type Encoder interface {
	Encode(id uint16, data Payload) ([]byte, error)
}

func NewParser(idToKind map[uint16]Kind, decoders map[Kind]Decoder, encoders map[Kind]Encoder) *Parser {
	return &Parser{
		idToKind: idToKind,
		decoders: decoders,
		encoders: encoders,
	}
}

func (parser *Parser) Decode(metadata Metadata, packet []byte) (Packet, error) {
	id, err := parser.getID(packet[:2])
	if err != nil {
		return Packet{}, err
	}
	metadata.ID = id

	decoder, err := parser.getDecoder(id)
	if err != nil {
		return Packet{}, err
	}

	payload, err := decoder.Decode(id, packet[2:])
	if err != nil {
		return Packet{}, err
	}

	return New(metadata, payload), nil
}

func (parser *Parser) getID(packet []byte) (uint16, error) {
	var id uint16
	err := binary.Read(bytes.NewReader(packet), binary.LittleEndian, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (parser *Parser) getDecoder(id uint16) (Decoder, error) {
	kind, ok := parser.idToKind[id]
	if !ok {
		return nil, fmt.Errorf("unknown packet %d", id)
	}

	decoder, ok := parser.decoders[kind]
	if !ok {
		return nil, fmt.Errorf("no decoder for packet %d", id)
	}

	return decoder, nil
}

func (parser *Parser) Encode(id uint16, payload Payload) ([]byte, error) {
	encoder, err := parser.getEncoder(id)
	if err != nil {
		return nil, err
	}

	idBytes := bytes.NewBuffer(nil)
	err = binary.Write(idBytes, binary.LittleEndian, id)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("idBytes: %v, id: %v", id, idBytes.Bytes())

	payloadBytes, err := encoder.Encode(id, payload)
	return append(idBytes.Bytes(), payloadBytes...), err
}

func (parser *Parser) getEncoder(id uint16) (Encoder, error) {
	kind, ok := parser.idToKind[id]
	if !ok {
		return nil, fmt.Errorf("unknown packet %d", id)
	}
	log.Info().Msgf("kind: %v", kind)

	encoder, ok := parser.encoders[kind]
	if !ok {
		return nil, fmt.Errorf("no encoder for packet %d", id)
	}

	return encoder, nil
}
