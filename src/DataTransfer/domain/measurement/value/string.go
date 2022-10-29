package value

import "fmt"

type String struct {
	value string
}

func (s *String) ToDisplayString() string {
	return fmt.Sprintf("%v", s.value)
}

func (s *String) Update(newValue any) {
	str, ok := newValue.(string)
	if !ok {
		panic("invalid value")
	}
	s.value = str
}
