package interfaces

type Structure interface {
	PacketName() string
	Measurements() []string
}
