package adapters

import (
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataStructure/infra/adapters/utils"
)

type PacketAdapter struct {
	Description  DescriptionAdapter
	Measurements []MeasurementAdapter
}

func expandPacket(description DescriptionAdapter, measurements []MeasurementAdapter) []PacketAdapter {
	ids := utils.GetAllIds(description.Id)
	packets := make([]PacketAdapter, 0)
	for index, id := range ids {
		newPacket := PacketAdapter{Description: description}
		newPacket.Description.Id = id
		sufix := "_" + strconv.Itoa(index)
		newPacket.Description.Name = newPacket.Description.Name + sufix
		newMeasurements := getMeasurementsWithSufix(sufix, measurements)
		newPacket.Measurements = newMeasurements
		packets = append(packets, newPacket)
	}

	return packets
}
