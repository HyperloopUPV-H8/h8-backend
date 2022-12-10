package models

import "time"

type Value struct {
	Name      string
	Value     any
	Timestamp time.Time
}
