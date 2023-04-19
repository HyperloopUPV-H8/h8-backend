package parsers

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
)

type parser interface {
	decode(packet.ValueDescriptor, binary.ByteOrder, io.Reader) (packet.Value, error)
	encode(packet.ValueDescriptor, binary.ByteOrder, packet.Value, io.Writer) error
}

type numericParser[T common.Numeric] struct{}

func (parser numericParser[T]) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	var value T
	err := binary.Read(data, order, &value)
	if err != nil {
		return packet.Numeric{}, err
	}

	return packet.Numeric{Value: (float64)(value)}, nil
}

func (parser numericParser[T]) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	return binary.Write(data, order, (T)(value.(packet.Numeric).Value))
}

type booleanParser struct{}

func (parser booleanParser) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	var value bool
	err := binary.Read(data, order, &value)
	if err != nil {
		return packet.Boolean{}, err
	}

	return packet.Boolean{Value: value}, nil
}

func (parser booleanParser) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	return binary.Write(data, order, value.(packet.Boolean).Value)
}

type enumParser struct {
	descriptors map[string]packet.EnumDescriptor
}

func (parser enumParser) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	enum, ok := parser.descriptors[descriptor.Name]
	if !ok {
		return packet.Enum{}, fmt.Errorf("enum descriptor for %s not found", descriptor.Name)
	}

	var value uint8
	err := binary.Read(data, order, &value)
	if err != nil {
		return packet.Enum{}, err
	}

	return packet.Enum{Value: enum[value]}, nil
}

func (parser enumParser) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	enum, ok := parser.descriptors[descriptor.Name]
	if !ok {
		return fmt.Errorf("enum descriptor for %s not found", descriptor.Name)
	}

	var index uint8
	for i, v := range enum {
		if v == value.(packet.Enum).Value {
			index = (uint8)(i)
			break
		}
	}

	return binary.Write(data, order, index)
}
