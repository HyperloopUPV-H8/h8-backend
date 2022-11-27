package domain

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain/idExpander"
)

type PacketDTO struct {
	Description  DescriptionDTO
	Measurements []MeasurementDTO
}

func expandPacket(description DescriptionDTO, measurements []MeasurementDTO) []PacketDTO {
	ids := idExpander.GetAllIds(description.ID)
	packets := make([]PacketDTO, len(ids))
	for index, id := range ids {
		newPacket := PacketDTO{Description: descriptionWithID(description, id), Measurements: measurements}
		if len(id) > 1 {
			sufix := fmt.Sprintf("_%v", index)
			newPacket = packetWithSufix(newPacket, sufix)
		}
		packets[index] = newPacket
	}

	return packets
}

func packetWithSufix(packet PacketDTO, sufix string) PacketDTO {
	return PacketDTO{
		Description:  descriptionWithName(packet.Description, packet.Description.Name+sufix),
		Measurements: measurementsWithSufix(packet.Measurements, sufix),
	}
}
