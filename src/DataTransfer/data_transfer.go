package application

import (
	"fmt"
	"os"

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

func (dataTransfer DataTransfer) Invoke(getPacketUpdate func() packetParser.PacketUpdate, logFile *os.File) {
	defer logFile.Close()
	for {
		update := getPacketUpdate()
		dataTransfer.data.UpdatePacket(update)
		packetTimestampPair := dataTransfer.data.GetPacket(update.ID)
		dataTransfer.PacketChannel <- packetTimestampPair.Packet

	}
}

func writePacket(packet *domain.Packet, logFile *os.File) {
	titlePacket := fmt.Sprintf(`Id: %v    Name: %v    Count: %v    CycleTime: %v`,
		packet.Id, packet.Name, packet.Count, packet.CycleTime)
	fmt.Fprintln(logFile, titlePacket)

	for _, measurement := range packet.Measurements {
		measuramentString := fmt.Sprintf(`	%v: %v`, measurement.Name, measurement.Value.ToDisplayString())
		fmt.Fprintln(logFile, measuramentString)
	}
}
