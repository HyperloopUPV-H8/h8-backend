package ade

import (
	"errors"

	doc "github.com/HyperloopUPV-H8/Backend-H8/excel/document"
)

const (
	Addresses  = "addresses"
	Units      = "units"
	Ports      = "ports"
	BoardIds   = "board_ids"
	MessageIds = "message_ids"
	BackendKey = "Backend"
)

func getInfo(sheet doc.Sheet) (Info, error) {
	tables, err := getTables(sheet)

	if err != nil {
		return Info{}, err
	}

	addresses, ok := tables[Addresses]

	if !ok {
		return Info{}, errors.New("addresses table not found")
	}

	addresses = removeHeaders(addresses)

	units, ok := tables[Units]

	if !ok {
		return Info{}, errors.New("units table not found")
	}

	units = removeHeaders(units)

	ports, ok := tables[Ports]

	if !ok {
		return Info{}, errors.New("ports table not found")
	}

	ports = removeHeaders(ports)

	boardIds, ok := tables[BoardIds]

	if !ok {
		return Info{}, errors.New("boardIds table not found")
	}

	boardIds = removeHeaders(boardIds)

	messageIds, ok := tables[MessageIds]

	if !ok {
		return Info{}, errors.New("messageIds table not found")
	}

	messageIds = removeHeaders(messageIds)

	return Info{
		Addresses:  tableToMap(addresses),
		Units:      tableToMap(units),
		Ports:      tableToMap(ports),
		BoardIds:   tableToMap(boardIds),
		MessageIds: tableToMap(messageIds),
	}, nil
}

func removeHeaders(table Table) Table {
	if len(table) == 0 {
		return table
	}

	return table[1:]
}

func tableToMap(table [][]string) map[string]string {
	m := make(map[string]string, len(table))

	for _, row := range table {
		m[row[0]] = row[1]
	}

	return m
}
