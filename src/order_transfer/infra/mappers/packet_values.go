package mappers

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer/application"
)

func GetPacketValues(order application.OrderJSON) dto.PacketValues {
	return dto.NewPacketValues(order.ID, order.Values)
}
