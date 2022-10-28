package dataTransfer

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	packetparserDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
	"github.com/davecgh/go-spew/spew"
)

type DataTransfer struct {
	pd domain.PodData
}

func New(podData domain.PodData) DataTransfer {
	dataTransfer := DataTransfer{
		pd: podData,
	}
	return dataTransfer
}

func (dt DataTransfer) Invoke(getPacketUpdate func() packetparserDomain.PacketUpdate) {

	for {
		update := getPacketUpdate()
		dt.pd.UpdatePacket(update)
		packet := dt.pd.GetPacket(update.ID)

		spew.Dump(packet)
		<-time.After(time.Second * 2)
	}
}
