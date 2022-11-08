package infra

import (
	"fmt"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
)

type DataTransfer struct {
	data                   domain.PodData
	PacketTimestampChannel chan domain.PacketTimestampPair
}

func New(rawBoards map[string]excelAdapter.BoardDTO) DataTransfer {
	return DataTransfer{
		data:                   domain.NewPodData(rawBoards),
		PacketTimestampChannel: make(chan domain.PacketTimestampPair),
	}
}

func (dataTransfer DataTransfer) Invoke(getPacketUpdate func() packetParser.PacketUpdate) {
	go func() {
		for {
			update := getPacketUpdate()
			fmt.Println(update)
			dataTransfer.data.UpdatePacket(update)
			packetTimestampPair := dataTransfer.data.GetPacket(update.ID)
			dataTransfer.PacketTimestampChannel <- *packetTimestampPair
		}
	}()
}
