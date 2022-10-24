package value

import (
	"fmt"
)

type Bool bool

// func (b *Bool) current() Value {
// 	return b
// }

// func (i *Bool) fromBuffer(b *bytes.Buffer) {
// 	binary.Read(b, binary.BigEndian, i)
// }

// func (b *Bool) toString() string {
// 	return fmt.Sprintf("%v", b)
// }

func (b *Bool) ToDisplayString() string {
	return fmt.Sprintf("%v", b)
}
