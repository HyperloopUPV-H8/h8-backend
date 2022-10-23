package value

import "bytes"

type UInt8 uint8
type UInt16 uint16
type UInt32 uint32
type UInt64 uint64

func (i *UInt8) current() Value {
	return i
}
func (i *UInt8) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 1)
	*i = UInt8(n)
}

func (i *UInt16) current() Value {
	return i
}
func (i *UInt16) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 2)
	*i = UInt16(n)
}

func (i *UInt32) current() Value {
	return i
}
func (i *UInt32) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 4)
	*i = UInt32(n)
}

func (i *UInt64) current() Value {
	return i
}
func (i *UInt64) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 8)
	*i = UInt64(n)
}
