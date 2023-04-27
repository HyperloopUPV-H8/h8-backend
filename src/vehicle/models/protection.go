package models

type ProtectionMessage struct {
	Kind       string     `json:"kind"`
	Board      string     `json:"board"`
	Name       string     `json:"name"`
	Protection Protection `json:"protection"`
	Timestamp  Timestamp  `json:"timestamp"`
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

type Error struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}
