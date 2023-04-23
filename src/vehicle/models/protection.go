package models

type Protection struct {
	Kind      string    `json:"kind"`
	Board     string    `json:"board"`
	Value     string    `json:"value"`
	Violation Violation `json:"violation"`
	Timestamp Timestamp `json:"timestamp"`
}

type Timestamp struct {
	Counter uint16 `json:"counter"`
	Seconds uint8  `json:"seconds"`
	Minutes uint8  `json:"minutes"`
	Hours   uint8  `json:"hours"`
	Day     uint8  `json:"day"`
	Month   uint8  `json:"month"`
	Year    uint16 `json:"year"`
}

type Violation interface {
	Type() string
}

type OutOfBoundsViolation struct {
	Kind string     `json:"kind"`
	Want [2]float64 `json:"want"`
	Got  float64    `json:"got"`
}

func (v OutOfBoundsViolation) Type() string {
	return v.Kind
}

type UpperBoundViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func (v UpperBoundViolation) Type() string {
	return v.Kind
}

type LowerBoundViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func (v LowerBoundViolation) Type() string {
	return v.Kind
}

type EqualsViolation struct {
	Kind string `json:"kind"`
	Got  float64
}

func (v EqualsViolation) Type() string {
	return v.Kind
}

type NotEqualsViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func (v NotEqualsViolation) Type() string {
	return v.Kind
}
