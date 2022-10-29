package board

import "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"

type Structure struct {
	packetName   string
	measurements []string
}

func (structure Structure) PacketName() string {
	return structure.packetName
}

func (structure Structure) Measurements() []string {
	return structure.measurements
}

func newStructure(column []string) interfaces.Structure {
	return Structure{
		packetName:   column[0],
		measurements: column[1:],
	}
}
