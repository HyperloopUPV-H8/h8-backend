package interfaces

import "time"

type PacketUpdate interface {
	ID() uint16
	Timestamp() time.Time
	Values() map[string]any
}
