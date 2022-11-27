package infra

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/domain"
)

type UnitAggregate struct {
	operations map[string]domain.Operations
}

func (aggregate UnitAggregate) Get(name string) domain.Operations {
	return aggregate.operations[name]
}

func NewPodUnitAggregate(boards map[string]excelAdapter.BoardDTO) *UnitAggregate {
	operations := make(map[string]domain.Operations)
	for _, board := range boards {
		for _, packet := range board.GetPackets() {
			for _, measure := range packet.Measurements {
				if measure.PodUnits != "" {
					operations[measure.Name] = domain.NewOperations(measure.PodUnits)
				}
			}
		}
	}

	return &UnitAggregate{
		operations: operations,
	}
}

func NewDisplayUnitAggregate(boards map[string]excelAdapter.BoardDTO) *UnitAggregate {
	operations := make(map[string]domain.Operations)
	for _, board := range boards {
		for _, packet := range board.GetPackets() {
			for _, measure := range packet.Measurements {
				if measure.DisplayUnits != "" {
					operations[measure.Name] = domain.NewOperations(measure.DisplayUnits)
				}
			}
		}
	}

	return &UnitAggregate{
		operations: operations,
	}
}
