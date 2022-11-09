package streaming

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/dto"
)

type PacketWebAdapter struct {
	Id                      uint16                  `json:"id"`
	Name                    string                  `json:"name"`
	MeasurementsWebAdapters []MeasurementWebAdapter `json:"measurements"`
	HexValue                string                  `json:"hexValue"`
	Count                   uint                    `json:"count"`
	CycleTime               uint                    `json:"cycleTime"`
}

func newPacketWebAdapter(packet dto.Packet) PacketWebAdapter {
	measurementWebAdapters := getMeasurementWebAdapters(packet.Measurements())
	return PacketWebAdapter{
		Id:                      packet.ID(),
		MeasurementsWebAdapters: measurementWebAdapters,
		HexValue:                fmt.Sprintf("%x", packet.HexValue()),
		Count:                   packet.Count(),
		CycleTime:               uint(packet.CycleTime().Milliseconds()),
	}
}

func getMeasurementWebAdapters(measurements map[string]any) []MeasurementWebAdapter {
	adapters := make([]MeasurementWebAdapter, 0, len(measurements))
	for name, value := range measurements {
		adapters = append(adapters, MeasurementWebAdapter{
			Name:  name,
			Value: fmt.Sprintf("%s", value),
		})
	}
	return adapters
}

type MeasurementWebAdapter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Units string `json:"units"`
}
