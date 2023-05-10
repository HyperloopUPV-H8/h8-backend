package message_parser

import (
	"encoding/json"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type MessageAdapter struct {
	BoardId    uint              `json:"boardId"`
	Timestamp  models.Timestamp  `json:"timestamp"`
	Protection ProtectionAdapter `json:"protection"`
}

type ProtectionAdapter struct {
	Name string           `json:"name"`
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
}
