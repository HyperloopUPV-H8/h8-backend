package pod_data

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/excel/ade"
	"github.com/HyperloopUPV-H8/Backend-H8/excel/utils"
)

const EnumType = "enum"

func getMeasurements(adeMeasurements []ade.Measurement, globalUnits map[string]utils.Operations) ([]Measurement, error) {
	measurements := make([]Measurement, 0)
	mErrors := common.NewErrorList()

	for _, adeMeas := range adeMeasurements {
		meas, err := getMeasurement(adeMeas, globalUnits)

		if err != nil {
			mErrors.Add(err)
			continue
		}

		measurements = append(measurements, meas)
	}

	if len(mErrors) > 0 {
		return nil, mErrors
	}

	return measurements, nil
}

func getMeasurement(adeMeas ade.Measurement, globalUnits map[string]utils.Operations) (Measurement, error) {
	if isNumeric(adeMeas.Type) {
		m, err := getNumericMeasurement(adeMeas, globalUnits)

		if err != nil {
			return nil, err
		}
		return m, nil
	} else if adeMeas.Type == "bool" {
		return getBooleanMeasurement(adeMeas), nil
	} else if strings.Contains(adeMeas.Type, "enum") {
		return getEnumMeasurement(adeMeas), nil
	} else if strings.Contains(adeMeas.Type, "array") {
		return getArrayMeasurement(adeMeas)
	} else {
		return nil, fmt.Errorf("invalid measurement type %s", adeMeas.Type)
	}
}

func getNumericMeasurement(adeMeas ade.Measurement, globalUnits map[string]utils.Operations) (NumericMeasurement, error) {
	measErrs := common.NewErrorList()

	safeRange, err := utils.ParseRange(adeMeas.SafeRange)

	if err != nil {
		measErrs.Add(err)
	}

	warningRange, err := utils.ParseRange(adeMeas.WarningRange)

	if err != nil {
		measErrs.Add(err)
	}

	displayUnits, err := utils.ParseUnits(adeMeas.DisplayUnits, globalUnits)

	if err != nil {
		measErrs.Add(err)
	}

	podUnits, err := utils.ParseUnits(adeMeas.PodUnits, globalUnits)

	if err != nil {
		measErrs.Add(err)
	}

	if len(measErrs) > 0 {
		return NumericMeasurement{}, measErrs
	}

	return NumericMeasurement{
		Id:           adeMeas.Id,
		Name:         adeMeas.Name,
		Type:         adeMeas.Type,
		Units:        displayUnits.Name,
		DisplayUnits: displayUnits,
		PodUnits:     podUnits,
		SafeRange:    safeRange,
		WarningRange: warningRange,
	}, nil
}

// TODO: THIS SHOULD ERROR
func getEnumMeasurement(adeMeas ade.Measurement) EnumMeasurement {
	return EnumMeasurement{
		Id:      adeMeas.Id,
		Name:    adeMeas.Name,
		Type:    EnumType,
		Options: getEnumMembers(adeMeas.Type),
	}
}

func getEnumMembers(enumExp string) []string {
	trimmedEnumExp := strings.Replace(enumExp, " ", "", -1)
	firstParenthesisIndex := strings.Index(trimmedEnumExp, "(")
	lastParenthesisIndex := strings.LastIndex(trimmedEnumExp, ")")

	return strings.Split(trimmedEnumExp[firstParenthesisIndex+1:lastParenthesisIndex], ",")
}

func getBooleanMeasurement(adeMeas ade.Measurement) BooleanMeasurement {
	return BooleanMeasurement{
		Id:   adeMeas.Id,
		Name: adeMeas.Name,
		Type: adeMeas.Type,
	}
}

var ArrayExp = regexp.MustCompile(`^array<(\w+)>$`)

func getArrayMeasurement(adeMeas ade.Measurement) (ArrayMeasurement, error) {
	submatches := ArrayExp.FindStringSubmatch(adeMeas.Type)

	if submatches == nil {
		return ArrayMeasurement{}, fmt.Errorf("invalid array type expression %s", adeMeas.Type)
	}

	if len(submatches) != 2 {
		return ArrayMeasurement{}, fmt.Errorf("invalid number of submatches %d", len(submatches))
	}

	itemType := submatches[1]

	var arr any
	switch itemType {
	case "uint8":
		arr = make([]uint8, 0)
	case "uint16":
		arr = make([]uint16, 0)
	case "uint32":
		arr = make([]uint32, 0)
	case "uint64":
		arr = make([]uint64, 0)
	case "int8":
		arr = make([]int8, 0)
	case "int16":
		arr = make([]int16, 0)
	case "int32":
		arr = make([]int32, 0)
	case "int64":
		arr = make([]int64, 0)
	case "float32":
		arr = make([]float32, 0)
	case "float64":
		arr = make([]float64, 0)
	case "bool":
		arr = make([]bool, 0)
	default:
		return ArrayMeasurement{}, fmt.Errorf("invalid item type %s", itemType)
	}

	return ArrayMeasurement{
		Id:       adeMeas.Id,
		Name:     adeMeas.Name,
		Type:     "array",
		ItemType: itemType,
		Value:    arr,
	}, nil
}

func isNumeric(kind string) bool {
	return (kind == "uint8" ||
		kind == "uint16" ||
		kind == "uint32" ||
		kind == "uint64" ||
		kind == "int8" ||
		kind == "int16" ||
		kind == "int32" ||
		kind == "int64" ||
		kind == "float32" ||
		kind == "float64")
}
