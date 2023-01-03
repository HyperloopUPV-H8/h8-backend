package models

type FromBoards interface {
	AddPacket(board string, ip string, desc Description, values []Value)
}
