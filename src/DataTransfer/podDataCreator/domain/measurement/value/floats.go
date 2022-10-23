package value

import "bytes"

type Float32 float32
type Float64 float64

func (f *Float32) current() Value {
	return f
}
func (i *Float32) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 4)
	*i = Float32(n)
}

func (f *Float64) current() Value {
	return f
}
func (i *Float64) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 8)
	*i = Float64(n)
}
