package interfaces

type Board interface {
	Name() string
	Descriptions() map[string]Description
	Measurements() map[string]Measurement
	Structure() map[string]Structure
	GetPackets() []Packet
}
