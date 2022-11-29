package infra

import (
	"testing"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
)

var glob any

func BenchmarkUnitsConvert(b *testing.B) {
	boards := getBoardsMock()

	podUnits := NewPodUnitAggregate(boards)
	displayUnits := NewDisplayUnitAggregate(boards)

	variables := []string{"Voltage0", "Voltage1", "Voltage2", "Voltage3", "Voltage4"}

	test := float64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test = Convert(test, podUnits.Get(variables[i%len(variables)]))
		test = Convert(test, displayUnits.Get(variables[i%len(variables)]))
	}
	glob = test
}

func BenchmarkUnitsRevert(b *testing.B) {
	boards := getBoardsMock()

	podUnits := NewPodUnitAggregate(boards)
	displayUnits := NewDisplayUnitAggregate(boards)

	variables := []string{"Voltage0", "Voltage1", "Voltage2", "Voltage3", "Voltage4"}

	test := float64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test = Revert(test, displayUnits.Get(variables[i%len(variables)]))
		test = Revert(test, podUnits.Get(variables[i%len(variables)]))
	}
	glob = test
}

func getBoardsMock() map[string]excelAdapter.BoardDTO {
	return map[string]excelAdapter.BoardDTO{
		"BMS": {
			Name: "BMS",
			Descriptions: map[string]excelAdapter.DescriptionDTO{
				"Voltages": {
					ID:        "0",
					Name:      "Voltages",
					Frecuency: "200",
					Direction: "Input",
					Protocol:  "UDP",
				},
			},
			Measurements: map[string]excelAdapter.MeasurementDTO{
				"Voltage0": {
					Name:         "Voltage0",
					ValueType:    "uint64",
					PodUnits:     "v#/100",
					DisplayUnits: "v#+100",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage1": {
					Name:         "Voltage1",
					ValueType:    "uint64",
					PodUnits:     "v#/100+100",
					DisplayUnits: "v#*100*5",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage2": {
					Name:         "Voltage2",
					ValueType:    "uint64",
					PodUnits:     "v#+200",
					DisplayUnits: "v#/5",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage3": {
					Name:         "Voltage3",
					ValueType:    "uint64",
					PodUnits:     "v#-100",
					DisplayUnits: "v#+42",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage4": {
					Name:         "Voltage4",
					ValueType:    "uint64",
					PodUnits:     "v#*5",
					DisplayUnits: "v#*25",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
			},
			Structures: map[string]excelAdapter.Structure{
				"Voltages": {
					PacketName:   "Voltages",
					Measurements: []string{"Voltage0", "Voltage1", "Voltage2", "Voltage3", "Voltage4"},
				},
			},
		},
	}
}
