package models

import "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"

type Board interface {
	Request(models.Order) error
	Listen(chan<- models.Update)
	Input(models.Update)
	Output(func(models.Order) error)
	Name() string
}
