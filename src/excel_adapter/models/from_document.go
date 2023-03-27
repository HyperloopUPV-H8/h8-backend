package models

type FromDocument interface {
	AddGlobal(GlobalInfo)
	AddPacket(boardName string, packet Packet)
}
