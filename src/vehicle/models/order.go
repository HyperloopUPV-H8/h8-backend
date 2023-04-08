package models

type Order struct {
	ID     uint16         `json:"id"`
	Fields map[string]any `json:"fields"`
}
