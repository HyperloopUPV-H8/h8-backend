package mappers

import (
	"log"
	"strconv"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func getPackets(packets []excelAdapter.PacketDTO) []domain.Packet {
	result := make([]domain.Packet, 0, len(packets))
	for _, packet := range packets {
		if packet.Description.Direction != "Input" {
			continue
		}
		result = append(result, getPacket(packet))
	}
	return result
}

func getPacket(packet excelAdapter.PacketDTO) domain.Packet {
	return domain.Packet{
		Id:           getID(packet.Description.ID),
		Name:         packet.Description.Name,
		HexValue:     "",
		Measurements: getMeasurements(packet.Measurements),
		Count:        0,
		CycleTime:    0,
	}
}

func getID(raw string) uint16 {
	id, err := strconv.ParseUint(raw, 10, 16)
	if err != nil {
		log.Fatalln(err)
	}

	return uint16(id)
}
