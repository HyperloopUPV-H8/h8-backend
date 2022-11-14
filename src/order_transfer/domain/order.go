package domain

type id = uint16

type Order struct {
	ID     id                `json:"id"`
	Values map[string]string `json:"values"`
}
