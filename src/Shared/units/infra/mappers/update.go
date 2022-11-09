package mappers

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra"
)

func ConvertUpdate(input dto.PacketUpdate, units infra.Units) dto.PacketUpdate {
	converted := make(map[string]any, len(input.Values()))
	for name, value := range input.Values() {
		converted[name] = units.ConvertDisplay(name, units.ConvertInternational(name, value))
	}
	return dto.NewPacketUpdate(input.ID(), converted, input.HexValue())
}
