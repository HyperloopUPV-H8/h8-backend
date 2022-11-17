package domain

type Board struct {
	Name    string   `json:"name"`
	Packets []Packet `json:"packets"`
}
