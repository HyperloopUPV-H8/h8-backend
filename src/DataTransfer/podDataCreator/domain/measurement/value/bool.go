package value

import (
	"fmt"
)

type Bool bool

func (b *Bool) ToDisplayString() string {
	return fmt.Sprintf("%v", b)
}
