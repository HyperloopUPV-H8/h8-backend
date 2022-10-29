package board

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board/idExpander"
)

type Packet struct {
	description  interfaces.Description
	measurements []interfaces.Measurement
}

func expandPacket(description interfaces.Description, measurements []interfaces.Measurement) []interfaces.Packet {
	ids := idExpander.GetAllIds(description.ID())
	packets := make([]interfaces.Packet, len(ids))
	for index, id := range ids {
		newPacket := Packet{description: descriptionWithID(description, id), measurements: measurements}
		sufix := fmt.Sprintf("_%v", index)
		newPacket = packetWithSufix(newPacket, sufix)
		packets[index] = newPacket
	}

	return packets
}

func packetWithSufix(packet Packet, sufix string) Packet {
	return Packet{
		description:  descriptionWithName(packet.description, packet.description.Name()+sufix),
		measurements: measurementsWithSufix(packet.measurements, sufix),
	}
}

func (packet Packet) Description() interfaces.Description {
	return packet.description
}

func (packet Packet) Measurements() []interfaces.Measurement {
	return packet.measurements
}
