package order_logger

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggableOrder models.PacketUpdate

func (lo LoggableOrder) Id() string {
	return fmt.Sprint(lo.Metadata.ID)
}

func (lo LoggableOrder) Log() []string {
	log := []string{fmt.Sprint(lo.Metadata.ID)}

	for name, value := range lo.Values {
		log = append(log, name, fmt.Sprint(value))
	}

	return log
}
