package message

import "strconv"

type Violation interface {
	Type() string
}

type OutOfBoundsViolation struct {
	Kind string     `json:"kind"`
	Want [2]float64 `json:"want"`
	Got  float64    `json:"got"`
}

func parseOutOfBounds(parts []string) (Violation, error) {
	violation := OutOfBoundsViolation{
		Kind: "OUT_OF_BOUNDS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want[0], err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	violation.Want[1], err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation OutOfBoundsViolation) Type() string {
	return violation.Kind
}

type UpperBoundViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func parseUpperBound(parts []string) (Violation, error) {
	violation := UpperBoundViolation{
		Kind: "UPPER_BOUND",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation UpperBoundViolation) Type() string {
	return "UPPER_BOUND"
}

type LowerBoundViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func parseLowerBound(parts []string) (Violation, error) {
	violation := LowerBoundViolation{
		Kind: "LOWER_BOUND",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation LowerBoundViolation) Type() string {
	return "LOWER_BOUND"
}

type EqualsViolation struct {
	// FIXME: do we need to parse the wanted value? is it included in the violation?
	Kind string `json:"kind"`
	Got  float64
}

func parseEquals(parts []string) (Violation, error) {
	violation := EqualsViolation{
		Kind: "EQUALS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation EqualsViolation) Type() string {
	return "EQUALS"
}

type NotEqualsViolation struct {
	Kind string `json:"kind"`
	Want float64
	Got  float64
}

func parseNotEquals(parts []string) (Violation, error) {
	violation := NotEqualsViolation{
		Kind: "NOT_EQUALS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation NotEqualsViolation) Type() string {
	return "NOT_EQUALS"
}

var violationStrategy = map[string]func([]string) (Violation, error){
	"OUT_OF_BOUNDS": parseOutOfBounds,
	"UPPER_BOUND":   parseUpperBound,
	"LOWER_BOUND":   parseLowerBound,
	"EQUALS":        parseEquals,
	"NOT_EQUALS":    parseNotEquals,
}

func parseViolation(data []string) (Violation, error) {
	return violationStrategy[data[0]](data[1:])
}
