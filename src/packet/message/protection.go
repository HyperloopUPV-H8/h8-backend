package message

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Protection struct {
	Kind      string
	Board     string
	Value     string
	Violation Violation
	Timestamp ProtectionTimestamp
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
	Counter uint16
	Seconds uint8
	Minutes uint8
	Hours   uint8
	Day     uint8
	Month   uint8
	Year    uint8
}

func parseTimestamp(data []byte) (ProtectionTimestamp, error) {
	var timestamp ProtectionTimestamp
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &timestamp)
	return timestamp, err
}
