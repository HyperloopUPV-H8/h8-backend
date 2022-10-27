package dto

import (
	"fmt"
	"log"
	"strconv"
	"time"

	podDataCreator "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto/idExpander"
)

type PacketDTO struct {
	Description  DescriptionDTO
	Measurements []MeasurementDTO
}

func (p PacketDTO) toPacket() podDataCreator.Packet {
	id, err := strconv.Atoi(p.Description.Id)

	if err != nil {
		log.Fatalf("parse: %s\n", err)
	}

	return podDataCreator.Packet{
		Id:           uint16(id),
		Name:         p.Description.Name,
		Measurements: p.getMeasurements(),
		Count:        0,
		CycleTime:    0,
		Timestamp:    time.Now(), //FIXME: que valor le pongo a esto?
	}
}

func (p PacketDTO) getMeasurements() map[string]*measurement.Measurement {
	measurements := make(map[Name]*measurement.Measurement, 0)
	for _, mDTO := range p.Measurements {
		measurement := mDTO.toMeasurement()
		measurements[measurement.Name] = &measurement
	}

	return measurements
}

func expandPacketDTO(description DescriptionDTO, measurements []MeasurementDTO) []PacketDTO {
	ids := idExpander.GetAllIds(description.Id)
	packetDTOs := make([]PacketDTO, len(ids))
	for index, id := range ids {
		newPacket := PacketDTO{Description: description, Measurements: measurements}
		newPacket.Description.Id = id
		sufix := fmt.Sprintf("_%v", index)
		newPacket = getWithSufix(newPacket, sufix)
		packetDTOs[index] = newPacket
	}

	return packetDTOs
}

func getWithSufix(packetDTO PacketDTO, sufix string) PacketDTO {
	packetDTO.Description.Name += sufix
	newMeasurements := getMeasurementsWithSufix(sufix, packetDTO.Measurements)
	packetDTO.Measurements = newMeasurements
	return packetDTO
}
