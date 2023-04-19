package models

type Update struct {
	ID        uint16                 `json:"id"`
	HexValue  string                 `json:"hexValue"`
	Count     uint64                 `json:"count"`
	CycleTime uint64                 `json:"cycleTime"`
	Values    map[string]UpdateValue `json:"measurementUpdates"`
}

type UpdateValue interface {
	Kind() string
}

type NumericValue struct {
	Value   float64 `json:"value"`
	Average float64 `json:"average"`
}

func (numeric NumericValue) Kind() string {
	return "numeric"
}

type BooleanValue struct {
	Value bool `json:"value"`
}

func (boolean BooleanValue) Kind() string {
	return "boolean"
}

type EnumValue struct {
	Value string `json:"value"`
}

func (enum EnumValue) Kind() string {
	return "enum"
}
