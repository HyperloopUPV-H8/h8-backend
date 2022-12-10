package models

type Order struct {
	ID     uint16
	Values map[string]any
}
