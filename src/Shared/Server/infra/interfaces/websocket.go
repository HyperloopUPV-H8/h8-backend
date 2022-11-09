package interfaces

type WebSocket interface {
	ReadJSON(any) error
	WriteJSON(any) error
	Close() error
}
