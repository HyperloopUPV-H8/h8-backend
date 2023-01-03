package infra

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/domain"
)

func GetMessage(update dto.PacketUpdate) domain.Message {
	return domain.NewMessage(update.ID(), update.Values())
}
