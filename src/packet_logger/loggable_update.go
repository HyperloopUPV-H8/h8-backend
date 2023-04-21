package packet_logger

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggablePacket struct {
	Metadata packet.Metadata
	HexValue []byte
}

func (packet LoggablePacket) Id() string {
	return strconv.Itoa(int(packet.Metadata.ID))
}

func (packet LoggablePacket) Log() []string {
	return []string{
		packet.Metadata.Timestamp.String(),
		packet.Metadata.From,
		packet.Metadata.To,
		fmt.Sprintf("%d", packet.Metadata.ID),
		fmt.Sprintf("%X", packet.HexValue),
	}
}

func ToLoggablePacket(update vehicle_models.PacketUpdate) LoggablePacket {
	return LoggablePacket{
		Metadata: update.Metadata,
		HexValue: update.HexValue,
	}
}

type LoggableValue struct {
	ValueId   string
	Value     packet.Value
	Timestamp time.Time // la del PacketUpdate
}

func (value LoggableValue) Id() string {
	return value.ValueId
}

func (value LoggableValue) Log() []string {
	return []string{
		value.ValueId,
		fmt.Sprintf("%v", value.Value),
		value.Timestamp.String(),
	}
}

func ToLoggableValue(id string, value packet.Value, timestamp time.Time) LoggableValue {
	return LoggableValue{
		ValueId:   id,
		Value:     value,
		Timestamp: timestamp,
	}
}
