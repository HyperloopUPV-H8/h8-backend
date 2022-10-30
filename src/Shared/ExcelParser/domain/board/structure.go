package board

type Structure struct {
	PacketName   string
	Measurements []string
}

func newStructure(column []string) Structure {
	return Structure{
		PacketName:   column[0],
		Measurements: column[1:],
	}
}
