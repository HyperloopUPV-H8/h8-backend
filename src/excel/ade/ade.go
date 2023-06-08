package ade

import (
	"github.com/HyperloopUPV-H8/Backend-H8/common"
	doc "github.com/HyperloopUPV-H8/Backend-H8/excel/document"
	"github.com/xuri/excelize/v2"
)

const (
	InfoName = "GLOBAL INFO" //TODO: change to INFO
)

func CreateADE(file *excelize.File) (ADE, error) {
	document := doc.CreateDocument(file)

	adeErrors := common.NewErrorList()

	info, err := getInfo(document)

	if err != nil {
		adeErrors.Add(err)
	}

	boardSheets := FilterMap(document.Sheets, func(name string, _ doc.Sheet) bool {
		return name != InfoName
	})

	boards, err := getBoards(boardSheets)

	if err != nil {
		adeErrors.Add(err)
	}

	return ADE{
		Info:   info,
		Boards: boards,
	}, adeErrors
}

func FilterMap[K comparable, V any](myMap map[K]V, predicate func(key K, value V) bool) map[K]V {
	filteredMap := make(map[K]V)

	for key, value := range myMap {
		if predicate(key, value) {
			filteredMap[key] = value
		}
	}

	return filteredMap
}
