package models

import "time"

type Value struct {
	Value     any
	Timestamp time.Time
}
