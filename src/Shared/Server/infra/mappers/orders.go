package mappers

import (
	"log"
	"strconv"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func GetOrders(boards map[string]excelAdapter.BoardDTO) map[uint16]domain.Order {
	orders := make(map[uint16]domain.Order)
	for _, board := range boards {
		for _, packet := range board.GetPackets() {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				log.Fatal(err)
			}

			fields := make([]domain.OrderField, len(packet.Measurements))
			for i, measurement := range packet.Measurements {
				fields[i] = domain.OrderField{
					Name:      measurement.Name,
					ValueType: measurement.ValueType,
				}
			}

			orders[uint16(id)] = domain.Order{
				ID:                uint16(id),
				Name:              packet.Description.Name,
				FieldDescriptions: fields,
			}
		}
	}
	return orders
}
