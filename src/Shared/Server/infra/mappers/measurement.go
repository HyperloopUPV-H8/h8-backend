package mappers

import (
	"strings"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func getMeasurements(measurements []excelAdapter.MeasurementDTO) []domain.Measurement {
	result := make([]domain.Measurement, 0, len(measurements))
	for _, measurement := range measurements {
		result = append(result, getMeasurement(measurement))
	}
	return result
}

func getMeasurement(measurement excelAdapter.MeasurementDTO) domain.Measurement {
	return domain.Measurement{
		Name:  measurement.Name,
		Value: "",
		Units: strings.Split(measurement.DisplayUnits, "#")[0],
		Type:  getType(measurement.ValueType),
	}
}

func getType(kind string) string {
	return "Number"
}
