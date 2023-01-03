package dto

import "time"

type LogValue struct {
	name      string
	timestamp time.Time
	data      string
}

func NewLogValue(name string, data string, timestamp time.Time) LogValue {
	return LogValue{
		name:      name,
		data:      data,
		timestamp: timestamp,
	}
}

func (value LogValue) Name() string {
	return value.name
}

func (value LogValue) Data() string {
	return value.data
}

func (value LogValue) Timestamp() time.Time {
	return value.timestamp
}
