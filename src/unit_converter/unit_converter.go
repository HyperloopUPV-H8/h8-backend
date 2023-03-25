package unit_converter

import (
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter/models"
)

type UnitConverter struct {
	Kind           string
	operations     map[string]models.Operations
	unitOperations map[string]string
}

func NewUnitConverter(kind string) *UnitConverter {
	return &UnitConverter{
		Kind:       kind,
		operations: make(map[string]models.Operations),
	}
}

func (converter *UnitConverter) AddGlobal(global excelAdapterModels.GlobalInfo) {
	converter.unitOperations = global.UnitToOperations
}

func (converter *UnitConverter) AddPacket(boardName string, packet excelAdapterModels.Packet) {
	for _, val := range packet.Values {
		if converter.Kind == "pod" && val.PodUnits != "" {
			converter.operations[val.ID] = models.NewOperations(getCustomOrGlobalOperations(val.PodUnits, converter.unitOperations))
		} else if converter.Kind == "display" && val.DisplayUnits != "" {
			converter.operations[val.ID] = models.NewOperations(getCustomOrGlobalOperations(val.DisplayUnits, converter.unitOperations))
		}
	}
}

func getCustomOrGlobalOperations(nameOrCustomOperations string, unitToOperations map[string]string) string {
	var operationsStr string

	if strings.Contains(nameOrCustomOperations, "#") {
		operationsStr = getCustomOperations(nameOrCustomOperations)
	} else {
		operationsStr = unitToOperations[nameOrCustomOperations]
	}

	return operationsStr
}

func getCustomOperations(operationsStr string) string {
	return strings.Split(operationsStr, "#")[1]
}

func (converter *UnitConverter) Convert(values map[string]any) map[string]any {
	convertedValues := make(map[string]any, len(values))
	for name, value := range values {
		if ops, ok := converter.operations[name]; ok {
			convertedValues[name] = ops.Convert(value.(float64))
		} else {
			convertedValues[name] = value
		}
	}
	return convertedValues
}

func (converter *UnitConverter) Revert(values map[string]any) map[string]any {
	convertedValues := make(map[string]any, len(values))
	for name, value := range values {
		if ops, ok := converter.operations[name]; ok {
			convertedValues[name] = ops.Convert(value.(float64))
		} else {
			convertedValues[name] = value
		}
	}
	return convertedValues
}
