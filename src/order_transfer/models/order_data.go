package models

import (
	"log"
	"strconv"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type OrderData map[string]OrderDescription

func (orderData *OrderData) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if orderData == nil {
		*orderData = make(OrderData)
	}

	id, err := strconv.ParseUint(desc.ID, 10, 16)
	if err != nil {
		log.Fatalf("order transfer: AddPacket: %s\n", err)
	}

	fields := make(map[string]FieldDescription, len(values))
	for _, value := range values {
		fields[value.Name] = FieldDescription{
			Type: value.Type,
		}
	}

	(*orderData)[desc.Name] = OrderDescription{
		ID:     uint16(id),
		Fields: fields,
	}
}

type OrderDescription struct {
	ID     uint16                      `json:"id"`
	Fields map[string]FieldDescription `json:"fieldDescriptions"`
}

type FieldDescription struct {
	Type string `json:"valueType"`
}
