package interfaces

type TransportController interface {
	ReceiveData() []byte
	ReceiveMessages() [][]byte
	Send(string, []byte)
	AliveConnections() []string
	Close()
}
