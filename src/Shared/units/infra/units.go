package infra

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/units/domain"
)

type Units map[string]domain.Operations

func NewUnits(literals map[string]string) Units {
	units := make(map[string]domain.Operations, len(literals))
	for name, literal := range literals {
		units[name] = domain.NewOperations(literal)
	}
	return units
}

func (units Units) Convert(name string, value any) any {
	if operations, exists := units[name]; exists {
		return domain.DoOperations(operations, value.(float64))
	}
	return value
}

func (units Units) Revert(name string, value any) any {
	if operations, exists := units[name]; exists {
		return domain.DoReverseOperations(operations, value.(float64))
	}
	return value
}
