package infra

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/units/domain"
)

type Units struct {
	pod     map[string]domain.Operations
	display map[string]domain.Operations
}

func NewUnits(podLiterals map[string]string, displayLiterals map[string]string) Units {
	return Units{
		pod:     getUnits(podLiterals),
		display: getUnits(displayLiterals),
	}
}

func getUnits(podLiterals map[string]string) map[string]domain.Operations {
	units := make(map[string]domain.Operations, len(podLiterals))
	for name, literal := range podLiterals {
		units[name] = domain.NewOperations(literal)
	}
	return units
}

func (units Units) ConvertInternational(name string, value any) any {
	if operations, exists := units.pod[name]; exists {
		return operations.Convert(value)
	}
	return value
}

func (units Units) ConvertDisplay(name string, value any) any {
	if operations, exists := units.display[name]; exists {
		return operations.Convert(value)
	}
	return value
}

func (units Units) RevertInternational(name string, value any) any {
	if operations, exists := units.pod[name]; exists {
		return operations.Revert(value)
	}
	return value
}

func (units Units) RevertPod(name string, value any) any {
	if operations, exists := units.display[name]; exists {
		return operations.Revert(value)
	}
	return value
}
