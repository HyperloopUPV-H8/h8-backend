package packet_adapter

import (
	"log"
	"testing"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra"
)

func BenchmarkParser(b *testing.B) {
	data := make(chan dto.PacketUpdate)
	adapter := New(infra.Config{
		Device:        "\\Device\\NPF_Loopback",
		Live:          true,
		TCPConfig:     nil,
		SnifferConfig: nil,
	}, 10, 0, data, nil, nil, getBoardsMock())

	b.ResetTimer()
	var last dto.PacketUpdate
	for i := 0; i < b.N; i++ {
		last = <-data
	}
	log.Println(last)
	log.Println(adapter)
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
				"Currents": {
					ID:        "1",
					Name:      "Currents",
					Frecuency: "200",
					Direction: "Input",
					Protocol:  "UDP",
				},
			},
			Measurements: map[string]excelAdapter.MeasurementDTO{
				"Voltage": {
					Name:         "Voltage",
					ValueType:    "uint8",
					PodUnits:     "v#",
					DisplayUnits: "v#",
					SafeRange:    "[110,120]",
					WarningRange: "[100,130]",
				},
				"Current": {
					Name:         "Current",
					ValueType:    "uint16",
					PodUnits:     "mA#/1000",
					DisplayUnits: "A#",
					SafeRange:    "[0,10]",
					WarningRange: "[-10,20]",
				},
			},
			Structures: map[string]excelAdapter.Structure{
				"Voltages": {
					PacketName:   "Voltages",
					Measurements: []string{"Voltage"},
				},
				"Currents": {
					PacketName:   "Currents",
					Measurements: []string{"Current"},
				},
			},
		},
	}
}
