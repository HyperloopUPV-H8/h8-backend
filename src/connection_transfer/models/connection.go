package models

type Connection struct {
	Name        string `json:"name"`
	IsConnected bool   `json:"isConnected"`
}
