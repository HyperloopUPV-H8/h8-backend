package board

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board/idExpander"
)

type Packet struct {
	Description  Description
	Measurements []Measurement
}

func expandPacket(description Description, measurements []Measurement) []Packet {
	ids := idExpander.GetAllIds(description.ID)
	packets := make([]Packet, len(ids))
	for index, id := range ids {
		newPacket := Packet{Description: descriptionWithID(description, id), Measurements: measurements}
		sufix := fmt.Sprintf("_%v", index)
		newPacket = packetWithSufix(newPacket, sufix)
		packets[index] = newPacket
	}

	return packets
}

func packetWithSufix(packet Packet, sufix string) Packet {
	return Packet{
		Description:  descriptionWithName(packet.Description, packet.Description.Name+sufix),
		Measurements: measurementsWithSufix(packet.Measurements, sufix),
	}
}
