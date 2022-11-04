package application

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain"
)

type DataTransfer struct {
	data          domain.PodData
	PacketChannel chan domain.Packet
}

func New(rawBoards map[string]excelParser.Board) DataTransfer {
	return DataTransfer{
		data:          domain.NewPodData(rawBoards),
		PacketChannel: make(chan domain.Packet, 10),
	}
}

func (dataTransfer DataTransfer) Invoke(getPacketUpdate func() packetParser.PacketUpdate) {
	for {
		update := getPacketUpdate()
		dataTransfer.data.UpdatePacket(update)
		packetTimestampPair := dataTransfer.data.GetPacket(update.ID)
		dataTransfer.PacketChannel <- packetTimestampPair.Packet
	}
}
