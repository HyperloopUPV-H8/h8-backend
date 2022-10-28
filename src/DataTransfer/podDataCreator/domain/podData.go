package podDataCreator

import (
	packetparser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packetParser"
)

type PodData struct {
	Boards map[string]*Board
}

func (podData *PodData) UpdatePacket(pu packetparser.PacketUpdate) {
	packet := podData.getPacket(pu.Id)
	packet.updatePacket(pu)
}

func (podData *PodData) getPacket(id uint16) *Packet {
	packet := new(Packet)
	for _, board := range podData.Boards {
		foundPacket, hasPacket := board.Packets[id]
		if hasPacket {
			packet = foundPacket
		}
	}
	return packet
}
