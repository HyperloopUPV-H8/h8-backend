package infra

import (
	"testing"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
)

var glob any

func BenchmarkParser(b *testing.B) {
	packetParser := NewPacketAggregate(getBoardsMock())

	var data dto.PacketUpdate
	for i := 0; i < b.N; i++ {
		data = Decode([]byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, *packetParser)
	}
	glob = data
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
					PodUnits:     "v#",
					DisplayUnits: "v#",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage1": {
					Name:         "Voltage1",
					ValueType:    "uint64",
					PodUnits:     "v#",
					DisplayUnits: "v#",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage2": {
					Name:         "Voltage2",
					ValueType:    "uint64",
					PodUnits:     "v#",
					DisplayUnits: "v#",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage3": {
					Name:         "Voltage3",
					ValueType:    "uint64",
					PodUnits:     "v#",
					DisplayUnits: "v#",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Voltage4": {
					Name:         "Voltage4",
					ValueType:    "uint64",
					PodUnits:     "v#",
					DisplayUnits: "v#",
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
