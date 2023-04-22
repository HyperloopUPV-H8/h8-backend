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
	log := []string{fmt.Sprint(lo.ID)}

	for name, field := range lo.Fields {
		log = append(log, name, fmt.Sprint(field.Value), fmt.Sprint(field.IsEnabled))
	}

	return log
}
