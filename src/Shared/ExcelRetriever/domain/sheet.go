package domain

type Sheet struct {
	Name   string
	Tables map[string]Table
}
