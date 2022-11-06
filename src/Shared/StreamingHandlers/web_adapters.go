package streaming

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type PacketWebAdapter struct {
	Id           uint16
	Name         string
	Measurements []MeasurementWebAdapter
	HexValue     []byte
	Count        uint
	CycleTime    uint
}

func newPacketWebAdapter(packet domain.Packet) PacketWebAdapter {
	measurementWebAdapters := getMeasurementWebAdapters(packet.Measurements)

	return PacketWebAdapter{
		Id:           packet.Id,
		Name:         packet.Name,
		Measurements: measurementWebAdapters,
		HexValue:     packet.HexValue,
		Count:        packet.Count,
		CycleTime:    uint(packet.CycleTime),
	}
}

func getMeasurementWebAdapters(measurements map[string]measurement.Measurement) []MeasurementWebAdapter {
	adapters := make([]MeasurementWebAdapter, len(measurements))
	for _, measurement := range measurements {
		adapters = append(adapters, MeasurementWebAdapter{
			Name:         measurement.Name,
			PodValue:     measurement.Value.ToPodUnitsString(),
			DisplayValue: measurement.Value.ToDisplayUnitsString(),
			PodUnits:     measurement.Value.GetPodUnits(),
			DisplayUnits: measurement.Value.GetDisplayUnits(),
		})
	}
	return adapters
}

type MeasurementWebAdapter struct {
	Name         string
	PodValue     string
	DisplayValue string
	PodUnits     string
	DisplayUnits string
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
