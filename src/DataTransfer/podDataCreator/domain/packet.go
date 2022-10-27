package domain

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
)

type Packet struct {
	Id           uint16
	Name         string
	Measurements map[string]*measurement.Measurement
	Count        uint
	CycleTime    int64
	Timestamp    time.Time
}

func (p *Packet) UpdatePacket(pu domain.PacketUpdate) {
	p.Count++
	p.CycleTime = pu.Timestamp.Sub(p.Timestamp).Milliseconds()
	p.Timestamp = pu.Timestamp
	for name, value := range pu.UpdatedValues {
		p.Measurements[name].Value.Update(value)
	}
}
