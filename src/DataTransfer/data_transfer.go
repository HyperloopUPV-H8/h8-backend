package application

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type DataTransfer struct {
	data          domain.PodData
	PacketChannel chan domain.Packet
}

func New(rawBoards map[string]excelAdapter.BoardDTO) DataTransfer {
	return DataTransfer{
		data:          domain.NewPodData(rawBoards),
		PacketChannel: make(chan domain.Packet, 10),
	}
}

func (dataTransfer DataTransfer) Parse(getPacketUpdate func() packetParser.PacketUpdate) *domain.PacketTimestampPair {
	update := getPacketUpdate()
	dataTransfer.data.UpdatePacket(update)
	return dataTransfer.data.GetPacket(update.ID)
}
