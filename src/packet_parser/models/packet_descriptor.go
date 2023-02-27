package models

type PacketDescriptor []ValueDescriptor

type ValueDescriptor struct {
	ID   string
	Type string
}

type PacketValues map[string]any
