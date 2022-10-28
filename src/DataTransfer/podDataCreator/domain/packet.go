package podDataCreator

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement"
	packetparser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packetParser"
)

type Packet struct {
	Id           uint16
	Name         string
	Measurements map[string]*measurement.Measurement
	Count        uint
	CycleTime    int64
	Timestamp    time.Time
}

func (p *Packet) updatePacket(pu packetparser.PacketUpdate) {
	p.Count++
	p.CycleTime = pu.Timestamp.Sub(p.Timestamp).Milliseconds()
	p.Timestamp = pu.Timestamp
	for name, value := range pu.MValues {
		p.Measurements[name].Value = value
	}
}
