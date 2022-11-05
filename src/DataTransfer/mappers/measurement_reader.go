package mappers

import (
	"fmt"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
)

type MeasurementReader struct {
	Measurement measurement.Measurement
	Timestamp   time.Time
}

func newMeausementReader(m measurement.Measurement, timestamp time.Time) MeasurementReader {
	return MeasurementReader{
		Measurement: m,
		Timestamp:   timestamp,
	}
}

func (mr MeasurementReader) Read(b []byte) (n int, err error) {
	b = []byte(fmt.Sprintf("%v", mr.Measurement.Value.ToDisplayUnitsString()))
	return len(b), nil
}

func GetMeasurementReaders(packetTimestampPair domain.PacketTimestampPair) []MeasurementReader {
	readers := make([]MeasurementReader, len(packetTimestampPair.Packet.Measurements))
	index := 0
	for _, measurement := range packetTimestampPair.Packet.Measurements {
		readers[index] = newMeausementReader(measurement, packetTimestampPair.Timestamp)
		index++
	}
	return readers
}
