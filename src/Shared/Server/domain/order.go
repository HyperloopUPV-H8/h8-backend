package domain

type Order struct {
	ID                uint16       `json:"id"`
	Name              string       `json:"name"`
	FieldDescriptions []OrderField `json:"fieldDescriptions"`
}

type OrderField struct {
	Name      string `json:"name"`
	ValueType string `json:"valueType"`
}
