package models

import (
	"fmt"
	"time"
)

type Value struct {
	Timestamp uint64
	Value     any
}

func NewValue(value any) Value {
	return Value{
		Timestamp: uint64(time.Now().UnixNano()),
		Value:     value,
	}
}

func (value *Value) ToCSV() []string {
	return []string{fmt.Sprint(value.Timestamp), fmt.Sprint(value.Value)}
}
