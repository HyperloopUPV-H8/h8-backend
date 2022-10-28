package dto

import excel "github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever"

type DescriptionDTO struct {
	Id        string
	Name      string
	Frecuency string
	Direction string
	Protocol  string
}

func newDescriptionDTO(row excel.Row) DescriptionDTO {
	return DescriptionDTO{
		Id:        row[0],
		Name:      row[1],
		Frecuency: row[2],
		Direction: row[3],
		Protocol:  row[4],
	}
}
