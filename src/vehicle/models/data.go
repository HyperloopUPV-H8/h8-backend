package models

type Data struct {
	ID        uint16         `json:"id"`
	HexValue  string         `json:"hexValue"`
	Count     uint64         `json:"count"`
	CycleTime uint64         `json:"cycleTime"`
	Fields    map[string]any `json:"measurementUpdates"`
}
