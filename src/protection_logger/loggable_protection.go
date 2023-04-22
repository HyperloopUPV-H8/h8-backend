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

	switch violation := lp.Violation.(type) {
	case vehicle_models.OutOfBoundsViolation:
		return []string{
			datetime,
			lp.Kind,
			lp.Board,
			lp.Value,
			violation.Kind,
			fmt.Sprintf("GOT %f", violation.Got),
			fmt.Sprintf("WANT %v", violation.Want),
		}
	case vehicle_models.LowerBoundViolation:
		return []string{
			datetime,
			lp.Kind,
			lp.Board,
			lp.Value,
			violation.Kind,
			fmt.Sprintf("GOT %f", violation.Got),
			fmt.Sprintf("WANT %f", violation.Want),
		}
	case vehicle_models.UpperBoundViolation:
		return []string{
			datetime,
			lp.Kind,
			lp.Board,
			lp.Value,
			violation.Kind,
			fmt.Sprintf("GOT %f", violation.Got),
			fmt.Sprintf("WANT < %f", violation.Want),
		}
	case vehicle_models.EqualsViolation:
		return []string{
			datetime,
			lp.Kind,
			lp.Board,
			lp.Value,
			violation.Kind,
			fmt.Sprintf("GOT %f", violation.Got),
		}
	case vehicle_models.NotEqualsViolation:
		return []string{
			datetime,
			lp.Kind,
			lp.Board,
			lp.Value,
			violation.Kind,
			fmt.Sprintf("GOT %f", violation.Got),
		}
	default:
		return []string{"MESSAGE NOT RECOGNIZED", lp.Board, lp.Kind, lp.Value}
	}

}
