package value

type Enum struct {
	values     map[uint]string
	currentVal uint
}

func (e *Enum) current() string {
	return e.values[e.currentVal]
}
