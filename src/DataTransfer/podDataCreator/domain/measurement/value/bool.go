package value

import "bytes"

type Bool bool

func (b *Bool) current() Value {
	return b
}
func (i *Bool) fromBuffer(b bytes.Buffer) {
	n := numberFromBuffer(b, 1)
	*i = n == 1
}
