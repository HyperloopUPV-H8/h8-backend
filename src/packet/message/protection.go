package message

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Protection struct {
	Kind      string              `json:"kind"`
	Board     string              `json:"board"`
	Value     string              `json:"value"`
	Violation Violation           `json:"violation"`
	Timestamp ProtectionTimestamp `json:"timestamp"`
	raw       []byte
}

func NewProtection(kind string, raw []byte) (Protection, error) {
	parts := strings.Split(string(raw), "\n")

	violation, err := parseViolation(parts[2:])
	if err != nil {
		return Protection{}, err
	}

	timestamp, err := parseTimestamp([]byte(parts[len(parts)-1]))
	if err != nil {
		return Protection{}, err
	}

	return Protection{
		Kind:      kind,
		Board:     parts[0],
		Value:     parts[1],
		Violation: violation,
		Timestamp: timestamp,
		raw:       raw,
	}, nil
}

func (message Protection) String() string {
	return string(message.raw)
}

type ProtectionTimestamp struct {
	Counter uint16 `json:"counter"`
	Seconds uint8  `json:"seconds"`
	Minutes uint8  `json:"minutes"`
	Hours   uint8  `json:"hours"`
	Day     uint8  `json:"day"`
	Month   uint8  `json:"month"`
	Year    uint8  `json:"year"`
}

func parseTimestamp(data []byte) (ProtectionTimestamp, error) {
	var timestamp ProtectionTimestamp
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &timestamp)
	return timestamp, err
}
