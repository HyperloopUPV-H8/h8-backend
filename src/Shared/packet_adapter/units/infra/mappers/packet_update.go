package mappers

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra"
)

func ConvertUpdate(update *dto.PacketUpdate, aggregate infra.UnitAggregate) {
	values := update.Values()

	for name, value := range values {
		if units := aggregate.Get(name); units != nil {
			update.SetValue(name, infra.Convert(value.(float64), units))
		}
	}
}

func RevertUpdate(update *dto.PacketUpdate, aggregate infra.UnitAggregate) {
	values := update.Values()

	for name, value := range values {
		if units := aggregate.Get(name); units != nil {
			update.SetValue(name, infra.Revert(value.(float64), units))
		}
	}
}
