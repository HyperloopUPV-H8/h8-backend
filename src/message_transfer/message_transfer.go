package messageTransfer

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
)

type MessageTransfer struct {
	MessageChannel chan dto.PacketUpdate
}

func New() MessageTransfer {
	return MessageTransfer{
		MessageChannel: make(chan dto.PacketUpdate),
	}
}

func (messageTransfer MessageTransfer) Invoke(getMessage func() dto.PacketUpdate) {
	for {
		messageTransfer.MessageChannel <- getMessage()
	}
}
