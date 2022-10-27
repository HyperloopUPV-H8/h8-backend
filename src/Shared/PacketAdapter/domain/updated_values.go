package domain

import (
	"time"

	value "github.com/HyperloopUPV-H8/Backend-H8/..."
)

type UpdatedValues struct {
	id        uint16
	measures  map[string]value.Value
	timestamp time.Time
}

func NewUpdatedValues(id ID, measures map[string]any) UpdatedValues {
	return UpdatedValues{
		id:        uint16(id),
		measures:  measures,
		timestamp: time.Now(),
	}
}
