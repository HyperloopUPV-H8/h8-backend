package infra

import "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/domain"

type UnitAggregate map[string]domain.Operations

func (aggregate UnitAggregate) Get(name string) domain.Operations {
	return aggregate[name]
}
