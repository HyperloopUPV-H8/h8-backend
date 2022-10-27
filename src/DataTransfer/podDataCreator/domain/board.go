package domain

type Board struct {
	Name    string
	Packets map[uint16]*Packet
}
