package packet

type Value any

type Numeric struct {
	Value float64
}

type Boolean struct {
	Value bool
}

type Enum struct {
	Value string
}
