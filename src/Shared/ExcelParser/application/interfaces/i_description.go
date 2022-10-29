package interfaces

type Description interface {
	ID() string
	Name() string
	Frecuency() string
	Direction() string
	Protocol() string
}
