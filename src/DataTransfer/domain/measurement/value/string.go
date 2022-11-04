package value

import "fmt"

type String string

func (s *String) ToPodUnitsString() string {
	return s.toString()
}

func (s *String) ToDisplayUnitsString() string {
	return s.toString()
}

func (s *String) GetPodUnits() string {
	return ""
}

func (s *String) GetDisplayUnits() string {
	return ""
}

func (s *String) toString() string {
	return fmt.Sprintf("%v", *s)
}

func (s *String) Update(newValue any) {
	str, ok := newValue.(string)
	if !ok {
		panic("invalid value")
	}
	*s = String(str)
}
