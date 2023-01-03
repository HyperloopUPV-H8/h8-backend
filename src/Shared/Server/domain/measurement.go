package domain

type Measurement struct {
	Name   string `json:"Name"`
	Value  string `json:"Value"`
	Ranges Ranges `json:"ranges"`
	Units  string `json:"units"`
}

type Ranges struct {
	Safe    [2]float64 `json:"safe"`
	Warning [2]float64 `json:"warning"`
}
