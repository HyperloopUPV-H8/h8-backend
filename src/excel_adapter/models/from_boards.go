package models

type FromBoards interface {
	AddPacket(globalInfo GlobalInfo, board string, ip string, desc Description, values []Value)
}
