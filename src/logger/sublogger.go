package logger

import "github.com/HyperloopUPV-H8/Backend-H8/packet"

type SubLogger interface {
	Start() error
	Stop() error
	Flush() error
	Close() error
	Update(packet.Packet) error
}
