package protection_logger

import (
	"fmt"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type LoggableProtection vehicle_models.Protection

func (lp LoggableProtection) Id() string {
	return lp.Kind
}

func (lp LoggableProtection) Log() []string {
	date := fmt.Sprintf("%d/%d/%d", lp.Timestamp.Day, lp.Timestamp.Month, lp.Timestamp.Year)
	time := fmt.Sprintf("%d:%d:%d", lp.Timestamp.Hours, lp.Timestamp.Minutes, lp.Timestamp.Seconds)
	datetime := fmt.Sprintf("%s %s", date, time)
	violationData := getViolationString(lp.Violation)
	return []string{datetime, lp.Kind, lp.Board, lp.Value, violationData}
}

func getViolationString(violation vehicle_models.Violation) string {
	var violationData string
	switch violation := violation.(type) {
	case vehicle_models.OutOfBoundsViolation:
		violationData = fmt.Sprintf("%s Got: %f Want: %f", violation.Kind, violation.Got, violation.Want)
	case vehicle_models.LowerBoundViolation:
		violationData = fmt.Sprintf("%s Got: %f Want: > %f", violation.Kind, violation.Got, violation.Want)
	case vehicle_models.UpperBoundViolation:
		violationData = fmt.Sprintf("%s Got: %f Want: < %f", violation.Kind, violation.Got, violation.Want)
	case vehicle_models.EqualsViolation:
		violationData = fmt.Sprintf("%s %f is not allowed", violation.Kind, violation.Got)
	case vehicle_models.NotEqualsViolation:
		violationData = fmt.Sprintf("%s %f should be %f", violation.Kind, violation.Got, violation.Want)
	default:
		violationData = fmt.Sprintf("%s UNRECOGNIZED VIOLATION", violation.Type())
	}

	return violationData
}
