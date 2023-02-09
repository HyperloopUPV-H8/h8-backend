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

func NewPodData() PodData {
	return PodData{
		Boards: make(map[string]Board),
	}
}

func (podData *PodData) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if desc.Type != "data" {
		return
	}

	id, err := strconv.ParseUint(desc.ID, 10, 16)
	if err != nil {
		log.Fatalf("data transfer: AddPacket: %s\n", err)
	}

	dataBoard, ok := podData.Boards[board]
	if !ok {
		podData.Boards[board] = Board{
			Name:    board,
			Packets: make(map[uint16]Packet),
		}
		dataBoard = podData.Boards[board]
	}

	dataBoard.Packets[uint16(id)] = Packet{
		ID:           uint16(id),
		Name:         desc.Name,
		HexValue:     "",
		Count:        0,
		CycleTime:    0,
		Measurements: getMeasurements(values),
	}
}

func getMeasurements(values []excelAdapterModels.Value) map[string]Measurement {
	measurements := make(map[string]Measurement, len(values))
	for _, value := range values {
		measurements[value.Name] = Measurement{
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

func parseRange(literal string) [2]float64 {
	if literal == "" {
		return [2]float64{0, 0}
	}

	split := strings.Split(strings.TrimSuffix(strings.TrimPrefix(literal, "["), "]"), ",")
	if len(split) != 2 {
		log.Fatalf("pod data: parseRange: invalid range %s\n", literal)
	}
	left, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		log.Fatalf("pod data: parseRange: %s\n", err)
	}

	right, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		log.Fatalf("pod data: parseRange: %s\n", err)
	}

	return [2]float64{left, right}
}

func getDefaultValue(valueType string) string {
	switch valueType {
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64":
		return "0"
	case "bool":
		return "false"
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
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Value        string     `json:"value"`
	Units        string     `json:"units"`
	SafeRange    [2]float64 `json:"safeRange"`
	WarningRange [2]float64 `json:"warningRange"`
}
