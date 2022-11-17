package application

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
)

type PacketJSON struct {
	ID        uint16            `json:"id"`
	HexValue  string            `json:"hexValue"`
	CycleTime uint64            `json:"cycleTime"`
	Count     uint              `json:"count"`
	Values    []MeasurementJSON `json:"measurementUpdates"`
}

type MeasurementJSON struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewJSON(packet domain.Packet) PacketJSON {
	return PacketJSON{
		ID:        packet.ID(),
		Count:     packet.Count(),
		CycleTime: uint64(packet.CycleTime().Milliseconds()),
		HexValue:  fmt.Sprintf("%x", packet.HexValue()),
		Values:    getValues(packet.Values()),
	}
}

func getValues(values map[string]any) []MeasurementJSON {
	result := make([]MeasurementJSON, 0, len(values))
	for name, value := range values {
		result = append(result, MeasurementJSON{
			Name:  name,
			Value: fmt.Sprintf("%v", value),
		})
	}
	return result
}
