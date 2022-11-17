package domain

type Packet struct {
	Id           uint16        `json:"id"`
	Name         string        `json:"name"`
	HexValue     string        `json:"hexValue"`
	Count        uint          `json:"count"`
	CycleTime    int64         `json:"cycleTime"`
	Measurements []Measurement `json:"measurements"`
}
