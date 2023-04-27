package protection_logger

import (
	"fmt"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggableProtection vehicle_models.ProtectionMessage

func (lp LoggableProtection) Id() string {
	return lp.Kind
}

func (lp LoggableProtection) Log() []string {
	date := fmt.Sprintf("%d/%d/%d", lp.Timestamp.Day, lp.Timestamp.Month, lp.Timestamp.Year)
	time := fmt.Sprintf("%d:%d:%d", lp.Timestamp.Hour, lp.Timestamp.Minute, lp.Timestamp.Second)
	datetime := fmt.Sprintf("%s %s", date, time)
	data := getDataString(lp.Protection.Data)
	return []string{datetime, lp.Kind, lp.Board, lp.Protection.Kind, data}
}

func getDataString(data any) string {
	switch data := data.(type) {
	case vehicle_models.OutOfBounds:
		return fmt.Sprintf("Got: %f Want: %f", data.Value, data.Bounds)
	case vehicle_models.LowerBound:
		return fmt.Sprintf("Got: %f Want: > %f", data.Value, data.Bound)
	case vehicle_models.UpperBound:
		return fmt.Sprintf("Got: %f Want: < %f", data.Value, data.Bound)
	case vehicle_models.Equals:
		return fmt.Sprintf("%f is not allowed", data.Value)
	case vehicle_models.NotEquals:
		return fmt.Sprintf("%f should be %f", data.Value, data.Want)
	default:
		return fmt.Sprintf("UNRECOGNIZED VIOLATION: %v", data)
	}

}
