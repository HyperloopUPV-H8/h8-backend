package domain

import (
	"log"
	"strconv"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type PacketTimestampPair struct {
	Packet    Packet
	Timestamp time.Time
}

type Packet struct {
	Id           uint16
	Name         string
	HexValue     []byte
	Measurements map[string]measurement.Measurement
	Count        uint
	CycleTime    int64
}

func (packetTimestampPair *PacketTimestampPair) UpdatePacket(pu packetParser.PacketUpdate) {
	packetTimestampPair.Packet.Count++
	packetTimestampPair.Packet.HexValue = pu.HexValue
	packetTimestampPair.Packet.CycleTime = pu.Timestamp.Sub(packetTimestampPair.Timestamp).Milliseconds()
	packetTimestampPair.Timestamp = pu.Timestamp
	for name, value := range pu.UpdatedValues {
		packetTimestampPair.Packet.Measurements[name].Value.Update(value)
	}
}

func NewPacketTimestampPairs(rawPackets []excelAdapter.PacketDTO) map[uint16]*PacketTimestampPair {
	packetTimestampPairs := make(map[uint16]*PacketTimestampPair, len(rawPackets))
	for _, packet := range rawPackets {
		id := getID(packet)
		packetTimestampPairs[id] = &PacketTimestampPair{
			Packet: Packet{
				Id:           id,
				Name:         packet.Description.Name,
				Measurements: measurement.NewMeasurements(packet.Measurements),
				Count:        0,
				CycleTime:    0,
			},
			Timestamp: time.Now(),
		}
	}
	return packetTimestampPairs
}

func getID(packet excelAdapter.PacketDTO) uint16 {
	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		log.Fatalf("get id: expected valid id, got %s\n", packet.Description.ID)
	}
	return uint16(id)
}
