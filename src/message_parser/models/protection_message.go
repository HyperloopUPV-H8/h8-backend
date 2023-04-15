package models

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	trace "github.com/rs/zerolog/log"
)

type ProtectionMessage struct {
	Kind      string              `json:"kind"`
	Board     string              `json:"board"`
	Value     string              `json:"value"`
	Violation any                 `json:"violation"`
	Timestamp ProtectionTimestamp `json:"timestamp"`
}

type ProtectionTimestamp struct {
	Counter uint16 `json:"counter"`
	Second  uint8  `json:"second"`
	Minute  uint8  `json:"minute"`
	Hour    uint8  `json:"hour"`
	Day     uint8  `json:"day"`
	Month   uint8  `json:"month"`
	Year    uint8  `json:"year"`
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

func ParseProtectionMessage(kind string, messageParts []string) (ProtectionMessage, error) {
	violation, err := getViolation(messageParts[2:])

	if err != nil {
		trace.Error().Err(err).Stack().Msg("parse violation")
		return ProtectionMessage{}, err
	}

	timestamp, err := getTimestamp(messageParts[len(messageParts)-1])

	if err != nil {
		trace.Error().Err(err).Stack().Msg("parse timestamp")
		return ProtectionMessage{}, err
	}

	return ProtectionMessage{
		Kind:      kind,
		Board:     messageParts[0],
		Value:     messageParts[1],
		Violation: violation,
		Timestamp: timestamp,
	}, nil
}

func getTimestamp(timestamp string) (ProtectionTimestamp, error) {
	var protectionTimestamp ProtectionTimestamp
	err := binary.Read(bytes.NewBuffer([]byte(timestamp)), binary.LittleEndian, &protectionTimestamp)
	return protectionTimestamp, err
}

func getViolation(violationParts []string) (any, error) {
	switch kind := violationParts[0]; kind {
	case "OUT_OF_BOUNDS":
		return getOutOfBoundsViolation(violationParts)
	case "UPPER_BOUND":
		return getUpperBoundViolation(violationParts)
	case "LOWER_BOUND":
		return getLowerBoundViolation(violationParts)
	case "EQUALS":
		return getEqualsViolation(violationParts)
	case "NOT_EQUALS":
		return getNotEqualsViolation(violationParts)
	default:
		return nil, errors.New("incorrect violation kind")
	}
}

func getOutOfBoundsViolation(violationParts []string) (OutOfBoundsViolation, error) {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		trace.Error().Err(gotErr).Stack().Msg("parse got")
		return OutOfBoundsViolation{}, gotErr
	}

	lowerBound, lowerErr := strconv.ParseFloat(violationParts[2], 64)

	if lowerErr != nil {
		trace.Error().Err(lowerErr).Stack().Msg("parse lower")
		return OutOfBoundsViolation{}, lowerErr
	}

	upperBound, upperErr := strconv.ParseFloat(violationParts[3], 64)

	if upperErr != nil {
		trace.Error().Err(upperErr).Stack().Msg("parse upper")
		return OutOfBoundsViolation{}, upperErr
	}

	return OutOfBoundsViolation{
		Want: [2]float64{lowerBound, upperBound},
		Got:  got,
	}, nil
}

func getUpperBoundViolation(violationParts []string) (UpperBoundViolation, error) {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		trace.Error().Err(gotErr).Stack().Msg("parse got")
		return UpperBoundViolation{}, gotErr
	}

	upperBound, upperErr := strconv.ParseFloat(violationParts[2], 64)

	if upperErr != nil {
		trace.Error().Err(upperErr).Stack().Msg("parse upper")
		return UpperBoundViolation{}, upperErr
	}

	return UpperBoundViolation{
		Want: upperBound,
		Got:  got,
	}, nil
}

func getLowerBoundViolation(violationParts []string) (LowerBoundViolation, error) {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		trace.Error().Err(gotErr).Stack().Msg("parse got")
		return LowerBoundViolation{}, gotErr
	}

	lowerBound, lowerErr := strconv.ParseFloat(violationParts[2], 64)

	if lowerErr != nil {
		trace.Error().Err(lowerErr).Stack().Msg("parse lower")
		return LowerBoundViolation{}, lowerErr
	}

	return LowerBoundViolation{
		Want: lowerBound,
		Got:  got,
	}, nil
}

func getEqualsViolation(violationParts []string) (EqualsViolation, error) {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		trace.Error().Err(gotErr).Stack().Msg("parse got")
		return EqualsViolation{}, gotErr
	}

	want, wantErr := strconv.ParseFloat(violationParts[2], 64)

	if wantErr != nil {
		trace.Error().Err(wantErr).Stack().Msg("parse want")
		return EqualsViolation{}, wantErr
	}

	return EqualsViolation{
		Want: want,
		Got:  got,
	}, nil
}

func getNotEqualsViolation(violationParts []string) (NotEqualsViolation, error) {
	got, gotErr := strconv.ParseFloat(violationParts[1], 64)

	if gotErr != nil {
		trace.Error().Err(gotErr).Stack().Msg("parse got")
		return NotEqualsViolation{}, gotErr
	}

	want, wantErr := strconv.ParseFloat(violationParts[2], 64)

	if wantErr != nil {
		trace.Error().Err(wantErr).Stack().Msg("parse want")
		return NotEqualsViolation{}, wantErr
	}

	return NotEqualsViolation{
		Want: want,
		Got:  got,
	}, nil
}
