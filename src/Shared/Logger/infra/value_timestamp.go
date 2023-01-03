package infra

import (
	"fmt"
	"time"
)

type ValueTimestamp struct {
	value     string
	timestamp time.Time
}

func NewValue(timestamp time.Time, value string) ValueTimestamp {
	return ValueTimestamp{
		value:     value,
		timestamp: timestamp,
	}
}

func (value ValueTimestamp) ToString() string {
	return fmt.Sprintf("%d, %s", value.timestamp.UnixNano(), value.value)
}
