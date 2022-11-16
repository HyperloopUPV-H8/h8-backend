package domain

type Board struct {
	Name    string            `json:"name"`
	Packets map[uint16]Packet `json:"packets"`
	Orders  map[uint16]Packet `json:"orders"`
}
