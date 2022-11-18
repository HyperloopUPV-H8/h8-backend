package application

type id = uint16

type OrderJSON struct {
	ID     id                `json:"id"`
	Values map[string]string `json:"values"`
}
