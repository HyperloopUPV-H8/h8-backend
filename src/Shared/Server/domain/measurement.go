package domain

type Measurement struct {
	Name  string `json:"Name"`
	Type  string `json:"type"`
	Value string `json:"Value"`
	Units string `json:"units"`
}
