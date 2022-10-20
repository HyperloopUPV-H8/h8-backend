package adapters

type StructureAdapter struct {
	packetName   string
	measurements []string
}

func newStructure(column []string) StructureAdapter {
	return StructureAdapter{
		packetName:   column[0],
		measurements: column[1:],
	}
}
