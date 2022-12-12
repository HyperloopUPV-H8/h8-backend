package unit_converter

import (
	"log"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter/models"
)

type UnitConverter struct {
	operations map[string]models.Operations
	Kind       string
}

func (converter *UnitConverter) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if converter.operations == nil {
		converter.operations = make(map[string]models.Operations)
	}

	for _, val := range values {
		if converter.Kind == "pod" {
			if val.PodOps == "" {
				continue
			}
			converter.operations[val.Name] = models.NewOperations(val.PodOps)
		} else if converter.Kind == "display" {
			if val.DisplayOps == "" {
				continue
			}
			converter.operations[val.Name] = models.NewOperations(val.DisplayOps)
		} else {
			log.Fatalf("unit converter: AddValue: invalid UnitConverter kind %s\n", converter.Kind)
		}
	}
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
