package datatransfer

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator"
	domain "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packetParser"
)

type DataTransfer struct {
	podData      domain.PodData
	packetParser packetParser.PacketParser
}

func (dataTransfer *DataTransfer) New() {
	dataTransfer.podData = podDataCreator.Run()
	dataTransfer.packetParser = packetParser.New()
}
