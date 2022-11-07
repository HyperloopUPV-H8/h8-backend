package value

import (
	"fmt"
)

type Bool bool

func (b Bool) ToDisplayUnitsString() string {
	return b.toString()
}

func (b *Bool) GetPodUnits() string {
	return ""
}

func (b *Bool) GetDisplayUnits() string {
	return ""
}

func (b Bool) ToPodUnitsString() string {
	return b.toString()
}

func (b Bool) toString() string {
	return fmt.Sprintf("%v", b)
}

func (b *Bool) Update(newValue any) {
	newBool, ok := newValue.(bool)
	if !ok {
		panic("invalid value")
	}
	*b = Bool(newBool)
}
