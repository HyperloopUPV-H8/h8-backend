package packetParser

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value"
)

type PacketParser struct {
}

type PacketUpdate struct {
	Id        uint16
	Timestamp time.Time
	MValues   map[string]value.Value
}

func New() PacketParser {
	return PacketParser{}
}
