package models

import "github.com/HyperloopUPV-H8/Backend-H8/packet"

type Order struct {
	ID     uint16           `json:"id"`
	Fields map[string]Field `json:"fields"`
}

type Field struct {
	Value     packet.Value `json:"value"`
	IsEnabled bool         `json:"isEnabled"`
}
