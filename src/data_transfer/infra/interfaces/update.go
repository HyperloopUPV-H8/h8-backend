package interfaces

import "time"

type id = uint16

type Update interface {
	ID() id
	Timestamp() time.Time
	HexValue() []byte
	Values() map[string]any
}
