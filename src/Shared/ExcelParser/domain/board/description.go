package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain"
)

type Description struct {
	id        string
	name      string
	frecuency string
	direction string
	protocol  string
}

func (desc Description) ID() string {
	return desc.id
}

func (desc Description) Name() string {
	return desc.name
}

func (desc Description) Frecuency() string {
	return desc.frecuency
}

func (desc Description) Direction() string {
	return desc.direction
}

func (desc Description) Protocol() string {
	return desc.protocol
}

func newDescription(row domain.Row) interfaces.Description {
	return Description{
		id:        row[0],
		name:      row[1],
		frecuency: row[2],
		direction: row[3],
		protocol:  row[4],
	}
}

func descriptionWithID(desc interfaces.Description, id string) interfaces.Description {
	return Description{
		id:        id,
		name:      desc.Name(),
		frecuency: desc.Frecuency(),
		direction: desc.Direction(),
		protocol:  desc.Protocol(),
	}
}

func descriptionWithName(desc interfaces.Description, name string) interfaces.Description {
	return Description{
		id:        desc.ID(),
		name:      name,
		frecuency: desc.Frecuency(),
		direction: desc.Direction(),
		protocol:  desc.Protocol(),
	}
}
