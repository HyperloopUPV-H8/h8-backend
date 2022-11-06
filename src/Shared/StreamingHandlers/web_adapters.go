package streaming

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type PacketWebAdapter struct {
	Id                      uint16                  `json:"id"`
	Name                    string                  `json:"name"`
	MeasurementsWebAdapters []MeasurementWebAdapter `json:"measurements"`
	HexValue                string                  `json:"hexValue"`
	Count                   uint                    `json:"count"`
	CycleTime               uint                    `json:"cycleTime"`
}

func newPacketWebAdapter(packet domain.Packet) PacketWebAdapter {
	measurementWebAdapters := getMeasurementWebAdapters(packet.Measurements)
	return PacketWebAdapter{
		Id:                      packet.Id,
		Name:                    packet.Name,
		MeasurementsWebAdapters: measurementWebAdapters,
		HexValue:                fmt.Sprintf("%x", packet.HexValue),
		Count:                   packet.Count,
		CycleTime:               uint(packet.CycleTime),
	}
}

func getMeasurementWebAdapters(measurements map[string]measurement.Measurement) []MeasurementWebAdapter {
	adapters := make([]MeasurementWebAdapter, len(measurements))
	index := 0
	for _, measurement := range measurements {
		adapters[index] = MeasurementWebAdapter{
			Name:         measurement.Name,
			PodValue:     measurement.Value.ToPodUnitsString(),
			DisplayValue: measurement.Value.ToDisplayUnitsString(),
			PodUnits:     measurement.Value.GetPodUnits(),
			DisplayUnits: measurement.Value.GetDisplayUnits(),
		}
		index++
	}
	return adapters
}

type MeasurementWebAdapter struct {
	Name         string `json:"name"`
	PodValue     string
	DisplayValue string `json:"value"`
	PodUnits     string
	DisplayUnits string `json:"units"`
}

type OrderWebAdapter struct {
	Id     uint16
	Fields map[string]string
}

type MessageWebAdapter struct {
	Id        uint16
	Fields    map[string]string
	Timestamp int64
}

func newMessageWebAdapter(packet packetParser.PacketUpdate) MessageWebAdapter {
	var value string
	var name string
	for n, val := range packet.UpdatedValues {
		value = val.(string)
		name = n
		break
	}

	return MessageWebAdapter{
		Id:        packet.ID,
		Fields:    map[string]string{name: value},
		Timestamp: packet.Timestamp.UnixNano(),
	}
}
