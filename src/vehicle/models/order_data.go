package models

import (
	"log"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type OrderData struct {
	Orders      map[string]OrderDescription `json:"orders"`
	StateOrders map[string]OrderDescription `json:"stateOrders"`
}

func NewOrderData(boards map[string]excelAdapterModels.Board, blcuName string) OrderData {
	orderData := OrderData{
		Orders:      make(map[string]OrderDescription),
		StateOrders: make(map[string]OrderDescription),
	}

	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type != "order" && packet.Description.Type != "stateOrder" {
				continue
			}

			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				log.Fatalf("order transfer: AddPacket: %s\n", err)
			}

			fields := make(map[string]FieldDescription, len(packet.Values))
			for _, value := range packet.Values {
				fields[value.Name] = getField(value.ID, value.Type, value.SafeRange, value.WarningRange)
			}

			desc := OrderDescription{
				ID:     uint16(id),
				Name:   packet.Description.Name,
				Fields: fields,
			}

			if packet.Description.Type == "order" {
				orderData.Orders[packet.Description.Name] = desc
			} else {
				orderData.StateOrders[packet.Description.Name] = desc
			}

		}
	}

	return orderData
}

func getField(name string, valueType string, safeRangeStr string, warningRangeStr string) FieldDescription {
	if IsNumeric(valueType) {

		SafeRange := parseRange(safeRangeStr)
		WarningRange := parseRange(warningRangeStr)

		return FieldDescription{
			Name: name,
			ValueDescription: NumericValue{
				Kind:         "numeric",
				Value:        valueType,
				SafeRange:    SafeRange,
				WarningRange: WarningRange,
			},
		}
	} else if valueType == "bool" {
		return FieldDescription{
			Name: name,
			ValueDescription: Value{
				Kind:  "boolean",
				Value: "",
			},
		}
	} else {
		return FieldDescription{
			Name: name,
			ValueDescription: Value{
				Kind:  "enum",
				Value: getEnumMembers(valueType),
			},
		}
	}
}

func getEnumMembers(enumExp string) []string {
	trimmedEnumExp := strings.Replace(enumExp, " ", "", -1)
	firstParenthesisIndex := strings.Index(trimmedEnumExp, "(")
	lastParenthesisIndex := strings.LastIndex(trimmedEnumExp, ")")

	return strings.Split(trimmedEnumExp[firstParenthesisIndex+1:lastParenthesisIndex], ",")
}

type OrderDescription struct {
	ID     uint16                      `json:"id"`
	Name   string                      `json:"name"`
	Fields map[string]FieldDescription `json:"fields"`
}

type FieldDescription struct {
	Name             string `json:"name"`
	ValueDescription any    `json:"valueDescription"`
}

type Value struct {
	Kind  string `json:"kind"`
	Value any    `json:"value"`
}

type NumericValue struct {
	Kind         string     `json:"kind"`
	Value        string     `json:"value"`
	SafeRange    []*float64 `json:"safeRange"`
	WarningRange []*float64 `json:"warningRange"`
}
