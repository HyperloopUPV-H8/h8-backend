package packet_adapter

import (
	"testing"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParserInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/sniffer"
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra"
	unitsMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra/mappers"
)

var glob any

func BenchmarkAdapter(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	sn := sniffer.New("\\Device\\NPF_Loopback", true, sniffer.DefaultConfig([]string{"127.0.0.2", "127.0.0.3"}, []string{"127.0.0.2", "127.0.0.3"}))

	pp := packetParserInfra.NewPacketAggregate(getBoardsMock())
	pu := unitsInfra.NewPodUnitAggregate(getBoardsMock())
	du := unitsInfra.NewDisplayUnitAggregate(getBoardsMock())

	var latest dto.PacketUpdate
	sn.GetNext()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		payload, _ := sn.GetNext()
		latest = packetParserInfra.Decode(payload, *pp)
		unitsMappers.ConvertUpdate(&latest, *pu)
		unitsMappers.ConvertUpdate(&latest, *du)
	}
	glob = latest
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
