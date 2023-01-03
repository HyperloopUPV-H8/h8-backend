package models

type Message struct {
	ID          uint16 `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
