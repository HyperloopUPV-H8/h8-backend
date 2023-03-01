package models

type PacketUpdate struct {
	ID        uint16         `json:"id"`
	HexValue  string         `json:"hexValue"`
	Count     uint64         `json:"count"`
	CycleTime uint64         `json:"cycleTime"`
	Values    map[string]any `json:"measurementUpdates"`
}
