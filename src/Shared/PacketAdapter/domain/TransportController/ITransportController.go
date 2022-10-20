package domain

type TransportController interface {
	Send(string, []byte)
	Receive() []byte
	Connected() []string
}
