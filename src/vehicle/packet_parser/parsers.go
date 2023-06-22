package packet_parser

import (
	"encoding/binary"
	"errors"
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
		return packet.Numeric(0), err
	}

	return packet.Numeric(value), nil
}

func (parser numericParser[T]) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	return binary.Write(data, order, (T)((value).(packet.Numeric)))
}

type booleanParser struct{}

func (parser booleanParser) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	var value bool
	err := binary.Read(data, order, &value)
	if err != nil {
		return packet.Boolean(false), err
	}

	return packet.Boolean(value), nil
}

func (parser booleanParser) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	return binary.Write(data, order, value)
}

type enumParser struct {
	descriptors map[string][]string
}

func (parser enumParser) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	enum, ok := parser.descriptors[descriptor.Name]
	if !ok {
		return packet.Enum("Default"), fmt.Errorf("enum descriptor for %s not found", descriptor.Name)
	}

	var value uint8
	err := binary.Read(data, order, &value)
	if err != nil {
		return packet.Enum("Default"), err
	}

	return packet.Enum(enum[value]), nil
}

func (parser enumParser) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	enum, ok := parser.descriptors[descriptor.Name]
	if !ok {
		return fmt.Errorf("enum descriptor for %s not found", descriptor.Name)
	}

	var index uint8
	for i, v := range enum {
		if v == string(value.(packet.Enum)) {
			index = (uint8)(i)
			break
		}
	}

	return binary.Write(data, order, index)
}

type arrayParser struct {
	idToItemType map[string]string
}

func (parser arrayParser) decode(descriptor packet.ValueDescriptor, order binary.ByteOrder, data io.Reader) (packet.Value, error) {
	var length uint8

	err := binary.Read(data, order, &length)

	if err != nil {
		return packet.Array{}, err
	}

	itemType, ok := parser.idToItemType[descriptor.Name]
	if !ok {
		return packet.Array{}, fmt.Errorf("enum descriptor for %s not found", descriptor.Name)
	}

	var arr any

	switch itemType {
	case "uint8":
		arr, err = readIntoArray[uint8](data, order, length)
	case "uint16":
		arr, err = readIntoArray[uint16](data, order, length)
	case "uint32":
		arr, err = readIntoArray[uint32](data, order, length)
	case "uint64":
		arr, err = readIntoArray[uint64](data, order, length)
	case "int8":
		arr, err = readIntoArray[int8](data, order, length)
	case "int16":
		arr, err = readIntoArray[int16](data, order, length)
	case "int32":
		arr, err = readIntoArray[int32](data, order, length)
	case "int64":
		arr, err = readIntoArray[int64](data, order, length)
	case "float32":
		arr, err = readIntoArray[float32](data, order, length)
	case "float64":
		arr, err = readIntoArray[float64](data, order, length)
	case "bool":
		arr, err = readIntoArray[bool](data, order, length)
	default:
		return packet.Array{}, err
	}

	if err != nil {
		return packet.Array{}, err
	}

	return packet.Array{
		Arr: arr,
	}, nil

}

func readIntoArray[T any](data io.Reader, order binary.ByteOrder, length uint8) (any, error) {
	arr := make([]T, length)
	err := binary.Read(data, order, arr)

	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (parser arrayParser) encode(descriptor packet.ValueDescriptor, order binary.ByteOrder, value packet.Value, data io.Writer) error {
	arr := value.Inner()

	switch typedArr := arr.(type) {
	case []uint8:
		binary.Write(data, order, typedArr)
	case []uint16:
		binary.Write(data, order, typedArr)
	case []uint32:
		binary.Write(data, order, typedArr)
	case []uint64:
		binary.Write(data, order, typedArr)
	case []int8:
		binary.Write(data, order, typedArr)
	case []int16:
		binary.Write(data, order, typedArr)
	case []int32:
		binary.Write(data, order, typedArr)
	case []int64:
		binary.Write(data, order, typedArr)
	case []float32:
		binary.Write(data, order, typedArr)
	case []float64:
		binary.Write(data, order, typedArr)
	case []bool:
		binary.Write(data, order, typedArr)
	default:
		return errors.New("invalid array type")
	}

	return nil
}
