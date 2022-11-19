package domain

type Measurement struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Units string `json:"units"`
}
