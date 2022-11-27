package mappers

import (
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer/application"
)

func GetPacketValues(order application.OrderJSON) dto.PacketUpdate {
	values := make(map[string]any, len(order.Values))
	for name, value := range order.Values {
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			values[name] = val
		} else if val, err := strconv.ParseBool(value); err == nil {
			values[name] = val
		} else {
			values[name] = val
		}
	}
	return dto.NewPacketUpdate(order.ID, values, []byte{})
}
