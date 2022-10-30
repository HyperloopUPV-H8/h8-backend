package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
)

type Description struct {
	ID        string
	Name      string
	Frecuency string
	Direction string
	Protocol  string
}

func newDescription(row document.Row) Description {
	return Description{
		ID:        row[0],
		Name:      row[1],
		Frecuency: row[2],
		Direction: row[3],
		Protocol:  row[4],
	}
}

func descriptionWithID(desc Description, id string) Description {
	return Description{
		ID:        id,
		Name:      desc.Name,
		Frecuency: desc.Frecuency,
		Direction: desc.Direction,
		Protocol:  desc.Protocol,
	}
}

func descriptionWithName(desc Description, name string) Description {
	return Description{
		ID:        desc.ID,
		Name:      name,
		Frecuency: desc.Frecuency,
		Direction: desc.Direction,
		Protocol:  desc.Protocol,
	}
}
