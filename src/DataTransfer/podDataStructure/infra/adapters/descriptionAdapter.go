package adapters

import excel "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelParser/domain"

type DescriptionAdapter struct {
	Id       string
	Name     string
	Frec     string
	Dir      string
	Protocol string
}

func newDescriptionAdapter(row excel.Row) DescriptionAdapter {
	return DescriptionAdapter{
		Id:       row[0],
		Name:     row[1],
		Frec:     row[2],
		Dir:      row[3],
		Protocol: row[4],
	}
}
