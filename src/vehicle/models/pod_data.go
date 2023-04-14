package models

import (
	"log"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type PodData struct {
	Boards map[string]Board `json:"boards"`
}

func NewPodData(excelBoards map[string]excelAdapterModels.Board) PodData {
	boards := make(map[string]Board)
	for name, excelBoard := range excelBoards {
		packets := make(map[uint16]Packet)
		for _, packet := range excelBoard.Packets {
			if packet.Description.Type != "data" {
				continue
			}

			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				log.Fatalf("data transfer: AddPacket: %s\n", err)
			}

			packets[uint16(id)] = Packet{
				ID:           uint16(id),
				Name:         packet.Description.Name,
				HexValue:     "",
				Count:        0,
				CycleTime:    0,
				Measurements: getMeasurements(packet.Values),
			}
		}

		boards[name] = Board{
			Name:    name,
			Packets: packets,
		}
	}

	return PodData{
		Boards: boards,
	}
}

func getMeasurements(values []excelAdapterModels.Value) map[string]Measurement {
	measurements := make(map[string]Measurement, len(values))
	for _, value := range values {
		measurements[value.ID] = Measurement{
			ID:   value.ID,
			Name: value.Name,
			Type: value.Type,
			//TODO: make sure added property (Value) doesn't break stuff
			Value:        getDefaultValue(value.Type),
			Units:        value.DisplayUnits,
			SafeRange:    parseRange(value.SafeRange),
			WarningRange: parseRange(value.WarningRange),
		}
	}
	return measurements
}

func parseRange(literal string) []*float64 {
	if literal == "" {
		return make([]*float64, 0)
	}

	strRange := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.Replace(literal, " ", "", -1), "["), "]"), ",")

	if len(strRange) != 2 {
		log.Fatalf("pod data: parseRange: invalid range %s\n", literal)
	}

	numRange := make([]*float64, 0)

	if strRange[0] != "" {
		lowerBound, errLowerBound := strconv.ParseFloat(strRange[0], 64)

		if errLowerBound != nil {
			log.Fatal("error parsing lower bound")
		}

		numRange = append(numRange, &lowerBound)
	} else {
		numRange = append(numRange, nil)
	}

	if strRange[1] != "" {
		upperBound, errUpperBound := strconv.ParseFloat(strRange[1], 64)

		if errUpperBound != nil {
			log.Fatal("error parsing lower bound")
		}

		numRange = append(numRange, &upperBound)
	} else {
		numRange = append(numRange, nil)
	}

	return numRange
}

func getDefaultValue(valueType string) any {
	switch valueType {
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64":
		return 0
	case "bool":
		return false
	default:
		return "Default"
	}
}

type Board struct {
	Name    string            `json:"name"`
	Packets map[uint16]Packet `json:"packets"`
}

type Packet struct {
	ID           uint16                 `json:"id"`
	Name         string                 `json:"name"`
	HexValue     string                 `json:"hexValue"`
	Count        uint16                 `json:"count"`
	CycleTime    int64                  `json:"cycleTime"`
	Measurements map[string]Measurement `json:"measurements"`
}

type Measurement struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Value        any        `json:"value"`
	Units        string     `json:"units"`
	SafeRange    []*float64 `json:"safeRange"`
	WarningRange []*float64 `json:"warningRange"`
}
