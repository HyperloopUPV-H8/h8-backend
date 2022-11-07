package messageTransfer

import packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"

type MessageTransfer struct {
	MessageChannel chan packetParser.PacketUpdate
}

func New() MessageTransfer {
	return MessageTransfer{
		MessageChannel: make(chan packetParser.PacketUpdate),
	}
}

func (messageTransfer MessageTransfer) Invoke(getMessage func() packetParser.PacketUpdate) {
	for {
		messageTransfer.MessageChannel <- getMessage()
	}
}
