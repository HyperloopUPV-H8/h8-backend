package unit_converter

import (
	"log"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter/models"
)

type UnitConverter struct {
	Kind       string
	operations map[string]models.Operations
}

func NewUnitConverter(kind string) UnitConverter {
	return UnitConverter{
		Kind:       kind,
		operations: make(map[string]models.Operations),
	}
}

func (converter *UnitConverter) AddPacket(globalInfo excelAdapterModels.GlobalInfo, board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	for _, val := range values {
		if converter.Kind == "pod" {
			converter.operations[val.Name] = models.NewOperations(getCustomOrGlobalOperations(val.PodUnits, globalInfo.UnitToOperations))
		} else if converter.Kind == "display" {
			converter.operations[val.Name] = models.NewOperations(getCustomOrGlobalOperations(val.DisplayUnits, globalInfo.UnitToOperations))
		} else {
			log.Fatalf("unit converter: AddValue: invalid UnitConverter kind %s\n", converter.Kind)
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
