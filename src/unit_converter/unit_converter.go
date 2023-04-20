package unit_converter

import (
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
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

func (converter *UnitConverter) Convert(values map[string]packet.Value) map[string]packet.Value {
	convertedValues := make(map[string]packet.Value, len(values))
	for name, originalValue := range values {
		numericValue, isNumeric := originalValue.(packet.Numeric)
		ops, ok := converter.operations[name]
		if ok && isNumeric {
			convertedValues[name] = packet.Numeric{Value: ops.Convert(numericValue.Value)}
		} else {
			convertedValues[name] = originalValue
		}
	}
	converter.trace.Trace().Msg("convert")
	return convertedValues
}

func (converter *UnitConverter) Revert(values map[string]packet.Value) map[string]packet.Value {
	convertedValues := make(map[string]packet.Value, len(values))
	for name, originalValue := range values {
		numericValue, isNumeric := originalValue.(packet.Numeric)
		ops, ok := converter.operations[name]
		if ok && isNumeric {
			convertedValues[name] = packet.Numeric{Value: ops.Revert(numericValue.Value)}
		} else {
			convertedValues[name] = originalValue
		}
	}
	converter.trace.Trace().Msg("convert")
	return convertedValues
}
