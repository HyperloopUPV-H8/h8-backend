package value

import "fmt"

type String string

func (s *String) ToDisplayString() string {
	return fmt.Sprintf("%v", *s)
}

func (s *String) Update(newValue any) {
	str, ok := newValue.(string)
	if !ok {
		panic("invalid value")
	}
	*s = String(str)
}
