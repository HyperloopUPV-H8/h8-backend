package domain

import "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"

type DescriptionDTO struct {
	ID        string
	Name      string
	Frecuency string
	Direction string
	Protocol  string
}

func newDescription(row domain.Row) DescriptionDTO {
	return DescriptionDTO{
		ID:        row[0],
		Name:      row[1],
		Frecuency: row[2],
		Direction: row[3],
		Protocol:  row[4],
	}
}

func descriptionWithID(desc DescriptionDTO, id string) DescriptionDTO {
	return DescriptionDTO{
		ID:        id,
		Name:      desc.Name,
		Frecuency: desc.Frecuency,
		Direction: desc.Direction,
		Protocol:  desc.Protocol,
	}
}

func descriptionWithName(desc DescriptionDTO, name string) DescriptionDTO {
	return DescriptionDTO{
		ID:        desc.ID,
		Name:      name,
		Frecuency: desc.Frecuency,
		Direction: desc.Direction,
		Protocol:  desc.Protocol,
	}
}
