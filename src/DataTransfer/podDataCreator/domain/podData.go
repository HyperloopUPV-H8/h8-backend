package domain

import (
	packetparser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
)

type PodData struct {
	Boards map[string]*Board
}

func (podData *PodData) UpdatePacket(pu packetparser.PacketUpdate) {
	packet := podData.getPacket(pu.ID)
	packet.UpdatePacket(pu)
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
