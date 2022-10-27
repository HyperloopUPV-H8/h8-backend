package podDataCreator

import (
	packetparser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packetParser"
)

type PodData struct {
	Boards map[string]*Board
}

func (podData *PodData) UpdatePacket(pu packetparser.PacketUpdate) {
	for _, board := range podData.Boards {
		packet, hasPacket := board.Packets[pu.Id]
		if hasPacket {
			packet.updatePacket(pu)
		}
	}

}
