package document

type Sheet struct {
	Name   string
	Tables map[string]Table
}
