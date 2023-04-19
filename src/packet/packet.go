package packet

type Packet struct {
	Metadata Metadata
	Payload  Payload
}

type Payload interface {
	Kind() Kind
}

type Kind int

const (
	Data Kind = iota
	Message
	Order
)

func New(metadata Metadata, payload Payload) Packet {
	return Packet{
		Metadata: metadata,
		Payload:  payload,
	}
}

func (packet Packet) Kind() Kind {
	return packet.Payload.Kind()
}
