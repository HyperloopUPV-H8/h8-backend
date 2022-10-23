package value

import "bytes"

type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64

func (i *Int8) current() Value {
	return i
}

// FIXME: el puntero paniquea si no ha sido inicializado -> crear con new()
func (i *Int8) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 1)
	*i = Int8(n)
}

func (i *Int16) current() Value {
	return i
}
func (i *Int16) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 2)
	*i = Int16(n)
}

func (i *Int32) current() Value {
	return i
}
func (i *Int32) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 4)
	*i = Int32(n)
}

func (i *Int64) current() Value {
	return i
}
func (i *Int64) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 8)
	*i = Int64(n)
}
