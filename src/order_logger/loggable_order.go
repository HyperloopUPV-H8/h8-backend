package order_logger

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggableOrder models.Order

func (lo LoggableOrder) Id() string {
	return fmt.Sprint(lo.ID)
}

func (lo LoggableOrder) Log() []string {
	return []string{fmt.Sprintf("[FROM GUI]: %d %v", lo.ID, lo.Fields)}
}

type LoggableTransmittedOrder models.PacketUpdate

func (lto LoggableTransmittedOrder) Id() string {
	return fmt.Sprint(lto.Metadata.ID)
}

func (lto LoggableTransmittedOrder) Log() []string {

	return []string{fmt.Sprintf("[TRANSMITTED]: %v %s %s %d %d %v", lto.Metadata.Timestamp, lto.Metadata.From, lto.Metadata.To, lto.Metadata.SeqNum, lto.Metadata.ID, lto.Values)}
}
