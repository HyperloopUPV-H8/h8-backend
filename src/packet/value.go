package packet

type Value interface {
	Inner() any
}

type Numeric struct {
	Value float64
}

func (n Numeric) Inner() any {
	return n.Value
}

type Boolean struct {
	Value bool
}

func (b Boolean) Inner() any {
	return b.Value
}

type Enum struct {
	Value string
}

func (e Enum) Inner() any {
	return e.Value
}
