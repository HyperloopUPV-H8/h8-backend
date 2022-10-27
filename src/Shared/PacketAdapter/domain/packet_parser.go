package domain

import "bytes"

type PacketParser struct {
	packets map[ID]Packet
	enums   map[Name]Enum
}

func (parser PacketParser) Decode(data []byte) PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := DecodeID(dataReader)

	values := parser.packets[id].Decode(parser.enums, dataReader)

	return NewPacketUpdate(id, values)
}
