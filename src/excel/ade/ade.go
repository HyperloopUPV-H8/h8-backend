package ade

import (
	"errors"

	doc "github.com/HyperloopUPV-H8/Backend-H8/excel/document"
	"github.com/xuri/excelize/v2"
)

const (
	InfoName = "GLOBAL INFO" //TODO: change to INFO
)

func CreateADE(file *excelize.File) (ADE, error) {
	document := doc.CreateDocument(file)
	infoSheet, ok := document.Sheets[InfoName]

	if !ok {
		return ADE{}, errors.New("info sheet not found")
	}

	info, err := getInfo(infoSheet)

	if err != nil {
		return ADE{}, err
	}

	boardSheets := FilterMap(document.Sheets, func(name string, _ doc.Sheet) bool {
		return name != InfoName
	})

	boards := getBoards(boardSheets)

	return ADE{
		Info:   info,
		Boards: boards,
	}, nil
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
