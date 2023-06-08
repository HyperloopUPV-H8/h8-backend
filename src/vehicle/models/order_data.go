package models

import (
	"log"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

const (
	OrderType      = "order"
	StateOrderType = "stateOrder"
	NumericKind    = "numeric"
	BooleanKind    = "boolean"
	EnumKind       = "enum"
)

type OrderData struct {
	Orders      map[string]OrderDescription `json:"orders"`
	StateOrders map[string]OrderDescription `json:"stateOrders"`
}

type VehicleOrders struct {
	Boards []BoardOrders `json:"boards"`
}

type BoardOrders struct {
	Name        string                  `json:"name"`
	Orders      []OrderDescription      `json:"orders"`
	StateOrders []StateOrderDescription `json:"stateOrders"`
}

type StateOrderDescription struct {
	OrderDescription
	Enabled bool `json:"enabled"`
}

func NewVehicleOrders(boards map[string]excelAdapterModels.Board, blcuName string) VehicleOrders {
	vehicleOrders := VehicleOrders{
		Boards: make([]BoardOrders, 0),
	}

	for _, board := range boards {
		boardOrders := BoardOrders{
			Name:        board.Name,
			Orders:      make([]OrderDescription, 0),
			StateOrders: make([]StateOrderDescription, 0),
		}

		for _, packet := range board.Packets {
			if packet.Description.Type != OrderType && packet.Description.Type != StateOrderType {
				continue
			}

			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				log.Fatalf("order transfer: AddPacket: %s\n", err)
			}

			fields := make(map[string]any, len(packet.Values))
			for _, value := range packet.Values {
				fields[value.Name] = getField(value.ID, value.Type, value.SafeRange, value.WarningRange)
			}

			desc := OrderDescription{
				ID:     uint16(id),
				Name:   packet.Description.Name,
				Fields: fields,
			}

			if packet.Description.Type == OrderType {
				boardOrders.Orders = append(boardOrders.Orders, desc)
			} else {
				boardOrders.StateOrders = append(boardOrders.StateOrders, StateOrderDescription{OrderDescription: desc, Enabled: false})
			}

		}
		vehicleOrders.Boards = append(vehicleOrders.Boards, boardOrders)
	}

	return vehicleOrders
}

func getField(name string, valueType string, safeRangeStr string, warningRangeStr string) any {
	if IsNumeric(valueType) {

		safeRange := parseRange(safeRangeStr)
		warningRange := parseRange(warningRangeStr)

		return NumericDescription{
			fieldDescription: fieldDescription{
				Kind: NumericKind,
				Name: name,
			},
			VarType:      valueType,
			SafeRange:    safeRange,
			WarningRange: warningRange,
		}
	} else if valueType == "bool" {
		return BooleanDescription{
			fieldDescription: fieldDescription{
				Kind: BooleanKind,
				Name: name,
			},
		}
	} else {
		return EnumDescription{
			fieldDescription: fieldDescription{
				Kind: EnumKind,
				Name: name,
			},
			Options: getEnumMembers(valueType),
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
	ID     uint16         `json:"id"`
	Name   string         `json:"name"`
	Fields map[string]any `json:"fields"`
}

type fieldDescription struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type NumericDescription struct {
	fieldDescription
	VarType      string     `json:"type"`
	SafeRange    []*float64 `json:"safeRange"`
	WarningRange []*float64 `json:"warningRange"`
}

type BooleanDescription struct {
	fieldDescription
}

type EnumDescription struct {
	fieldDescription
	Options []string `json:"options"`
}
