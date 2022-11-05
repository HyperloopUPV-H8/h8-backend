package domain

import "io"

type Entry struct {
	Id    string
	Value io.Reader
}
