package message

import "strconv"

type Violation interface {
	Kind() string
}

type OutOfBoundsViolation struct {
	Want [2]float64
	Got  float64
}

func parseOutOfBounds(parts []string) (Violation, error) {
	var violation OutOfBoundsViolation
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

func (violation OutOfBoundsViolation) Kind() string {
	return "OUT_OF_BOUNDS"
}

type UpperBoundViolation struct {
	Want float64
	Got  float64
}

func parseUpperBound(parts []string) (Violation, error) {
	var violation UpperBoundViolation
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

func (violation UpperBoundViolation) Kind() string {
	return "UPPER_BOUND"
}

type LowerBoundViolation struct {
	Want float64
	Got  float64
}

func parseLowerBound(parts []string) (Violation, error) {
	var violation LowerBoundViolation
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

func (violation LowerBoundViolation) Kind() string {
	return "LOWER_BOUND"
}

type EqualsViolation struct {
	// FIXME: do we need to parse the wanted value? is it included in the violation?
	Got float64
}

func parseEquals(parts []string) (Violation, error) {
	var violation EqualsViolation
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func (violation EqualsViolation) Kind() string {
	return "EQUALS"
}

type NotEqualsViolation struct {
	Want float64
	Got  float64
}

func parseNotEquals(parts []string) (Violation, error) {
	var violation NotEqualsViolation
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

func (violation NotEqualsViolation) Kind() string {
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
