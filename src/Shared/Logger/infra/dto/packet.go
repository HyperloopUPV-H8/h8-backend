package dto

import (
	"time"
)

type LogPacket struct {
	timestamp time.Time
	values    []LogValue
}

func NewLogPacket(timestamp time.Time, values []LogValue) LogPacket {
	return LogPacket{
		timestamp: timestamp,
		values:    values,
	}
}

func (packet LogPacket) Values() []LogValue {
	return packet.values
}

func (packet LogPacket) Timestamp() time.Time {
	return packet.timestamp
}
