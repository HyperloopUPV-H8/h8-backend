package packet

type Raw struct {
	Metadata Metadata
	Payload  []byte
}
