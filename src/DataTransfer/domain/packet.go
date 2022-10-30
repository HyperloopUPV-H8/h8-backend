package domain

import (
	"log"
	"strconv"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain"
)

type Packet struct {
	Id           uint16
	Name         string
	Measurements map[string]measurement.Measurement
	Count        uint
	CycleTime    int64
	Timestamp    time.Time
}

func (packet *Packet) UpdatePacket(data packetParser.PacketUpdate) {
	packet.Count++
	packet.CycleTime = data.Timestamp.Sub(packet.Timestamp).Milliseconds()
	packet.Timestamp = data.Timestamp
	for name, value := range data.UpdatedValues {
		packet.Measurements[name].Value.Update(value)
	}
}

func NewPackets(rawPackets []excelParser.Packet) map[uint16]Packet {
	packets := make(map[uint16]Packet, len(rawPackets))
	for _, packet := range rawPackets {
		id := getID(packet)
		packets[id] = Packet{
			Id:           id,
			Name:         packet.Description.Name,
			Measurements: measurement.NewMeasurements(packet.Measurements),
			Count:        0,
			CycleTime:    0,
			Timestamp:    time.Now(),
		}
	}
	return packets
}

func getID(packet excelParser.Packet) uint16 {
	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		log.Fatalf("get id: expected valid id, got %s\n", packet.Description.ID)
	}
	return uint16(id)
}
