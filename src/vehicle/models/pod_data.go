package models

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type PodData struct {
	Boards []Board `json:"boards"`
}

func NewPodData(excelBoards map[string]excelAdapterModels.Board) PodData {
	boards := make([]Board, 0)
	for name, excelBoard := range excelBoards {
		boards = append(boards, Board{
			Name:    name,
			Packets: getPackets(excelBoard.Packets),
		})
	}

	return PodData{
		Boards: boards,
	}
}

func getPackets(excelPackets []excelAdapterModels.Packet) []Packet {
	packets := make([]Packet, 0)

	for _, packet := range excelPackets {

		packet, err := getPacket(packet)

		if err != nil {
			continue
		}

		packets = append(packets, packet)
	}

	sortedPackets := SortablePacket(packets)
	sort.Sort(sortedPackets)

	return sortedPackets
}

func getPacket(packet excelAdapterModels.Packet) (Packet, error) {
	if packet.Description.Type != "data" {
		return Packet{}, fmt.Errorf("packet %s is not data packet", packet.Description.Name)
	}

	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		log.Fatalf("data transfer: AddPacket: %s\n", err)
	}

	return Packet{
		ID:           uint16(id),
		Name:         packet.Description.Name,
		HexValue:     "000000",
		Count:        0,
		CycleTime:    0,
		Measurements: getMeasurements(packet.Values),
	}, nil
}

func getMeasurements(values []excelAdapterModels.Value) []any {
	measurements := make([]any, 0)
	for _, value := range values {
		if IsNumeric(value.Type) {
			measurements = append(measurements, getNumericMeasurement(value))
		} else if value.Type == "bool" {
			measurements = append(measurements, getBooleanMeasurement(value))
		} else {
			measurements = append(measurements, getEnumMeasurement(value))
		}
	}
	return measurements
}

func getNumericMeasurement(value excelAdapterModels.Value) NumericMeasurement {
	return NumericMeasurement{
		ID:           value.ID,
		Name:         value.Name,
		Type:         value.Type,
		Units:        value.DisplayUnits,
		SafeRange:    parseRange(value.SafeRange),
		WarningRange: parseRange(value.WarningRange),
	}
}

func getEnumMeasurement(value excelAdapterModels.Value) EnumMeasurement {
	return EnumMeasurement{
		ID:   value.ID,
		Name: value.Name,
		Type: "Enum",
	}
}

func getBooleanMeasurement(value excelAdapterModels.Value) BooleanMeasurement {
	return BooleanMeasurement{
		ID:   value.ID,
		Name: value.Name,
		Type: value.Type,
	}
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

type Board struct {
	Name    string   `json:"name"`
	Packets []Packet `json:"packets"`
}

type Packet struct {
	ID           uint16 `json:"id"`
	Name         string `json:"name"`
	HexValue     string `json:"hexValue"`
	Count        uint16 `json:"count"`
	CycleTime    int64  `json:"cycleTime"`
	Measurements []any  `json:"measurements"`
}

type NumericMeasurement struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Units        string     `json:"units"`
	SafeRange    []*float64 `json:"safeRange"`
	WarningRange []*float64 `json:"warningRange"`
}

type BooleanMeasurement struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type EnumMeasurement struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
