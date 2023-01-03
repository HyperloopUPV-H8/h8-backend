package dto

type PacketValues struct {
	id     id
	values map[string]string
}

func NewPacketValues(id id, values map[string]string) PacketValues {
	return PacketValues{
		id:     id,
		values: values,
	}
}

func (values PacketValues) ID() id {
	return values.id
}

func (values PacketValues) GetValue(name string) string {
	return values.values[name]
}
