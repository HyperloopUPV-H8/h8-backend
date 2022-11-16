package domain

type Packet struct {
	Id           uint16                 `json:"id"`
	Name         string                 `json:"name"`
	Frequency    uint                   `json:"frequency"`
	Direction    string                 `json:"direction"`
	Protocol     string                 `json:"direction"`
	HexValue     string                 `json:"hexValue"`
	Measurements map[string]Measurement `json:"measurements"`
	Count        uint                   `json:"count"`
	CycleTime    int64                  `json:"cycleTime"`
}
