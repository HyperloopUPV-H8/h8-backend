package infra

import (
	"fmt"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement/value"
)

type ValueTimestamp struct {
	value     string
	timestamp time.Time
}

func NewValue(timestamp time.Time, measurement value.Value) ValueTimestamp {
	return ValueTimestamp{
		value:     measurement.ToDisplayString(),
		timestamp: timestamp,
	}
}

func (value ValueTimestamp) ToString() string {
	return fmt.Sprintf("%d, %s", value.timestamp.UnixNano(), value.value)
}
