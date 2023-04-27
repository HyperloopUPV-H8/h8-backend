package models

type ErrorMessage struct {
	Board      string          `json:"board"`
	Timestamp  Timestamp       `json:"timestamp"`
	Protection ErrorProtection `json:"protection"`
}

type ErrorProtection struct {
	Name string    `json:"name"`
	Type string    `json:"type"`
	Data ErrorData `json:"data"`
}

type ErrorData struct {
	Value string `json:"value"`
}
