package application

import (
	"fmt"
	"os"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain"
)

type DataTransfer struct {
	data domain.PodData
}

func New(rawBoards map[string]excelParser.Board) DataTransfer {
	return DataTransfer{
		data: domain.NewPodData(rawBoards),
	}
}

func (dataTransfer DataTransfer) Invoke(getPacketUpdate func() packetParser.PacketUpdate, logFile *os.File) {

	for {
		update := getPacketUpdate()
		dataTransfer.data.UpdatePacket(update)
		packet := dataTransfer.data.GetPacket(update.ID)

		titlePacket := fmt.Sprintf(`Id: %v    Name: %v    Count: %v    CycleTime: %v`,
			packet.Id, packet.Name, packet.Count, packet.CycleTime)
		fmt.Fprintln(logFile, titlePacket)

		for _, measurement := range packet.Measurements {
			measuramentString := fmt.Sprintf(`	%v: %v`, measurement.Name, measurement.Value.ToDisplayString())
			fmt.Fprintln(logFile, measuramentString)
		}
	}
}
