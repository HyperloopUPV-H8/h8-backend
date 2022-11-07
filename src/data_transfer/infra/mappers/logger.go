package mappers

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
)

func ToLogPacket(packet domain.PacketTimestampPair) dto.LogPacket {
	values := make([]dto.LogValue, 0, len(packet.Packet.Measurements))
	for name, measure := range packet.Packet.Measurements {
		values = append(values, dto.NewLogValue(name, measure.Value.GetDisplayUnits()))
	}
	return dto.NewLogPacket(packet.Timestamp, values)
}
