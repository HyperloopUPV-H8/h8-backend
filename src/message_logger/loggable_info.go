package protection_logger

import (
	"fmt"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggableInfo vehicle_models.InfoMessage

func (lp LoggableInfo) Id() string {
	return "info"
}

func (lp LoggableInfo) Log() []string {
	date := fmt.Sprintf("%d/%d/%d", lp.Timestamp.Day, lp.Timestamp.Month, lp.Timestamp.Year)
	time := fmt.Sprintf("%d:%d:%d", lp.Timestamp.Hour, lp.Timestamp.Minute, lp.Timestamp.Second)
	datetime := fmt.Sprintf("%s %s", date, time)
	return []string{datetime, "info", lp.Board, lp.Msg}
}
