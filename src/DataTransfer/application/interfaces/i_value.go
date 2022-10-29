package interfaces

type Value interface {
	ToDisplayString() string
	Update(any)
}
