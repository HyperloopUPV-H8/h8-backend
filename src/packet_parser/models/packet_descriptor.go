package models

type PacketDescriptor []ValueDescriptor

type ValueDescriptor struct {
	Name string
	Type string
}

type PacketValues map[string]any
