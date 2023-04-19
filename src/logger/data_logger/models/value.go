package models

import (
	"fmt"
	"time"

	update_factory_models "github.com/HyperloopUPV-H8/Backend-H8/update_factory/models"
)

type Value struct {
	Timestamp uint64
	Value     any
}

func NewValue(value update_factory_models.UpdateValue) Value {
	var data any
	switch value := value.(type) {
	case update_factory_models.NumericValue:
		data = value.Value
	case update_factory_models.EnumValue:
		data = value.Value
	case update_factory_models.BooleanValue:
		data = value.Value
	}

	return Value{
		Timestamp: uint64(time.Now().UnixNano()),
		Value:     data,
	}
}

func (value *Value) ToCSV() []string {
	return []string{fmt.Sprint(value.Timestamp), fmt.Sprint(value.Value)}
}
