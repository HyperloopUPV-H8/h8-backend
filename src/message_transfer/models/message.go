package models

import (
	"errors"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

const (
	FAULT_FIELD   = "fault"
	FAULT_TYPE    = "fault"
	WARNING_FIELD = "warning"
	WARNING_TYPE  = "warning"
)

type Message struct {
	ID          uint16 `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func MessageFromUpdate(update vehicle_models.Update) (Message, error) {
	if fault, ok := update.Fields[FAULT_FIELD]; ok {
		return Message{
			ID:          update.ID,
			Description: fault.(string),
			Type:        FAULT_TYPE,
		}, nil
	}

	if warning, ok := update.Fields[WARNING_FIELD]; ok {
		return Message{
			ID:          update.ID,
			Description: warning.(string),
			Type:        WARNING_FIELD,
		}, nil
	}

	return Message{}, errors.New("update isn't message")
}
