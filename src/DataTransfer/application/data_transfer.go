package application

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/interfaces"

	"github.com/davecgh/go-spew/spew"
)

type DataTransfer struct {
	data domain.PodData
}

func New(rawBoards map[string]interfaces.Board) DataTransfer {
	return DataTransfer{
		data: domain.NewPodData(rawBoards),
	}
}

func (dataTransfer DataTransfer) Invoke(getPacketUpdate func() packetAdapter.PacketUpdate) {
	for {
		update := getPacketUpdate()
		dataTransfer.data.UpdatePacket(update)
		packet := dataTransfer.data.GetPacket(update.ID())

		spew.Dump(packet)
		<-time.After(time.Second)
	}
}
