package application

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
)

type PacketJSON struct {
	ID        uint16            `json:"id"`
	Count     uint              `json:"count"`
	CycleTime uint64            `json:"cycleTime"`
	HexValue  string            `json:"hexValue"`
	Values    map[string]string `json:"values"`
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

func getValues(values map[string]any) map[string]string {
	result := make(map[string]string, len(values))
	for name, value := range values {
		result[name] = fmt.Sprintf("%v", value)
	}
	return result
}
