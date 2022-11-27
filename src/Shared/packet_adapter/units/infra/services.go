package infra

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/domain"
)

func Convert(value float64, operations domain.Operations) float64 {
	return domain.DoOperations(operations, value)
}

func Revert(value float64, operations domain.Operations) float64 {
	return domain.DoReverseOperations(operations, value)
}
