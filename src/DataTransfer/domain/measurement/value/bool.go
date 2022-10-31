package value

import (
	"fmt"
)

type Bool bool

func (b Bool) ToDisplayString() string {
	return fmt.Sprintf("%v", b)
}

func (b *Bool) Update(newValue any) {
	newBool, ok := newValue.(bool)
	if !ok {
		panic("invalid value")
	}
	*b = Bool(newBool)
}
