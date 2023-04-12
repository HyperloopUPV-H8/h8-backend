package unit_converter

import (
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type UnitConverter struct {
	operations map[string]models.Operations
	trace      zerolog.Logger
}

func NewUnitConverter(kind string, boards map[string]excelAdapterModels.Board, unitToOperations map[string]string) UnitConverter {
	trace.Info().Str("kind", kind).Msg("new unit converter")

	operations := make(map[string]models.Operations)

	for _, board := range boards {
		for _, packet := range board.Packets {
			for _, val := range packet.Values {
				var ops string
				if kind == "pod" && val.PodUnits != "" {
					ops = getCustomOrGlobalOperations(val.PodUnits, unitToOperations)
				} else if kind == "display" && val.DisplayUnits != "" {
					ops = getCustomOrGlobalOperations(val.DisplayUnits, unitToOperations)
				} else {
					continue
				}
				operations[val.ID] = models.NewOperations(ops)
			}
		}
	}

	return UnitConverter{
		operations: operations,
		trace:      trace.With().Str("component", "unitConverter").Str("kind", kind).Logger(),
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
	converter.trace.Trace().Msg("convert")
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
	converter.trace.Trace().Msg("convert")
	return convertedValues
}
