package mappers

import (
	"log"
	"strconv"
	"strings"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func getMeasurements(measurements []excelAdapter.MeasurementDTO) map[string]domain.Measurement {
	result := make(map[string]domain.Measurement, len(measurements))
	for _, measurement := range measurements {
		result[measurement.Name] = getMeasurement(measurement)
	}
	return result
}

func getMeasurement(measurement excelAdapter.MeasurementDTO) domain.Measurement {
	return domain.Measurement{
		Name:   measurement.Name,
		Value:  "",
		Units:  strings.Split(measurement.DisplayUnits, "#")[0],
		Ranges: getRanges(measurement.SafeRange, measurement.WarningRange),
	}
}

func getRanges(safe, warning string) domain.Ranges {
	return domain.Ranges{
		Safe:    getRange(safe),
		Warning: getRange(warning),
	}
}

func getRange(literal string) [2]float64 {
	parts := strings.Split(strings.Trim(literal, "[ ]"), ",")
	start, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 64)
	if err != nil {
		log.Fatalln(err)
	}
	end, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 64)
	if err != nil {
		log.Fatalln(err)
	}

	return [2]float64{start, end}
}
