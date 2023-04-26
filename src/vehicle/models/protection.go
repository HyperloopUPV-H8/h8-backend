package models

type ProtectionMessage struct {
	Kind       string     `json:"kind"`
	Board      string     `json:"board"`
	Name       string     `json:"name"`
	Protection Protection `json:"protection"`
	Timestamp  Timestamp  `json:"timestamp"`
}

type Timestamp struct {
	Counter uint16 `json:"counter"`
	Second  uint8  `json:"second"`
	Minute  uint8  `json:"minute"`
	Hour    uint8  `json:"hour"`
	Day     uint8  `json:"day"`
	Month   uint8  `json:"month"`
	Year    uint16 `json:"year"`
}

type Protection struct {
	Kind string `json:"kind"`
	Data any    `json:"data"`
}

type OutOfBounds struct {
	Value  float64    `json:"value"`
	Bounds [2]float64 `json:"bounds"`
}

type LowerBound struct {
	Value float64 `json:"value"`
	Bound float64 `json:"bound"`
}

type UpperBound struct {
	Value float64 `json:"value"`
	Bound float64 `json:"bound"`
}

type Equals struct {
	Value float64 `json:"value"`
}

type NotEquals struct {
	Value float64 `json:"value"`
	Want  float64 `json:"want"`
}
