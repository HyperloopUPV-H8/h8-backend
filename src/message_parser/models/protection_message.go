package models

import (
	"errors"
	"strconv"
)

type ProtectionMessage struct {
	Kind      string `json:"kind"`
	Board     string `json:"board"`
	Value     string `json:"value"`
	Violation any    `json:"violation"`
}

type OutOfBoundsViolation struct {
	Want [2]float64 `json:"want"`
	Got  float64    `json:"got"`
}

type UpperBoundViolation struct {
	Want float64 `json:"want"`
	Got  float64 `json:"got"`
}

type LowerBoundViolation struct {
	Want float64 `json:"want"`
	Got  float64 `json:"got"`
}

type EqualsViolation struct {
	Want float64 `json:"want"`
	Got  float64 `json:"got"`
}

type NotEqualsViolation struct {
	Want float64 `json:"want"`
	Got  float64 `json:"got"`
}

func ParseProtectionMessage(kind string, messageParts []string) ProtectionMessage {
	violation, err := getViolation(messageParts[2:])

	if err != nil {
		//TODO: error
	}

	return ProtectionMessage{
		Kind:      kind,
		Board:     messageParts[0],
		Value:     messageParts[1],
		Violation: violation,
	}
}

func getViolation(violationParts []string) (any, error) {
	switch kind := violationParts[0]; kind {
	case "OUT_OF_BOUNDS":
		return getOutOfBoundsViolation(violationParts), nil
	case "UPPER_BOUND":
		return getUpperBoundViolation(violationParts), nil
	case "LOWER_BOUND":
		return getLowerBoundViolation(violationParts), nil
	case "EQUALS":
		return getEqualsViolation(violationParts), nil
	case "NOT_EQUALS":
		return getNotEqualsViolation(violationParts), nil
	default: //TODO: undo kind of message
		return nil, errors.New("incorrect violation kind")
	}
}

func getOutOfBoundsViolation(violationParts []string) OutOfBoundsViolation {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		//TODO: error

	}

	lowerBound, lowerErr := strconv.ParseFloat(violationParts[2], 64)

	if lowerErr != nil {
		//TODO: error
	}

	upperBound, upperErr := strconv.ParseFloat(violationParts[3], 64)

	if upperErr != nil {
		//TODO: error
	}

	return OutOfBoundsViolation{
		Want: [2]float64{lowerBound, upperBound},
		Got:  got,
	}
}

func getUpperBoundViolation(violationParts []string) UpperBoundViolation {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		//TODO: error

	}

	upperBound, upperErr := strconv.ParseFloat(violationParts[2], 64)

	if upperErr != nil {
		//TODO: error

	}

	return UpperBoundViolation{
		Want: upperBound,
		Got:  got,
	}
}

func getLowerBoundViolation(violationParts []string) LowerBoundViolation {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		//TODO: error

	}

	lowerBound, lowerErr := strconv.ParseFloat(violationParts[2], 64)

	if lowerErr != nil {
		//TODO: error
	}

	return LowerBoundViolation{
		Want: lowerBound,
		Got:  got,
	}
}

func getEqualsViolation(violationParts []string) EqualsViolation {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		//TODO: error

	}

	want, wantErr := strconv.ParseFloat(violationParts[2], 64)

	if wantErr != nil {
		//TODO: error
	}

	return EqualsViolation{
		Want: want,
		Got:  got,
	}
}

func getNotEqualsViolation(violationParts []string) NotEqualsViolation {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		//TODO: error
	}

	want, wantErr := strconv.ParseFloat(violationParts[2], 64)

	if wantErr != nil {
		//TODO: error
	}

	return NotEqualsViolation{
		Want: want,
		Got:  got,
	}
}
