package domain

type id = uint16

type Message struct {
	id     id
	values map[string]any
}

func NewMessage(id id, values map[string]any) Message {
	return Message{
		id:     id,
		values: values,
	}
}

func (msg Message) ID() id {
	return msg.id
}

func (msg Message) Values() map[string]any {
	return msg.values
}
