package unit_converter

import (
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type UnitConverter struct {
	Kind           string
	operations     map[string]models.Operations
	unitOperations map[string]string
	trace          zerolog.Logger
}

func NewUnitConverter(kind string) *UnitConverter {
	trace.Info().Str("kind", kind).Msg("new unit converter")
	return &UnitConverter{
		Kind:       kind,
		operations: make(map[string]models.Operations),
		trace:      trace.With().Str("component", "unitConverter").Str("kind", kind).Logger(),
	}
}

func (converter *UnitConverter) AddGlobal(global excelAdapterModels.GlobalInfo) {
	converter.trace.Debug().Msg("add global")
	converter.unitOperations = global.UnitToOperations
}

func (converter *UnitConverter) AddPacket(boardName string, packet excelAdapterModels.Packet) {
	converter.trace.Debug().Str("id", packet.Description.ID).Str("name", packet.Description.Name).Str("board", boardName).Msg("add packet")
	for _, val := range packet.Values {
		var ops string
		if converter.Kind == "pod" && val.PodUnits != "" {
			ops = getCustomOrGlobalOperations(val.PodUnits, converter.unitOperations)
		} else if converter.Kind == "display" && val.DisplayUnits != "" {
			ops = getCustomOrGlobalOperations(val.DisplayUnits, converter.unitOperations)
		} else {
			converter.trace.Trace().Str("id", val.ID).Msg("no unit conversion")
			continue
		}
		converter.trace.Trace().Str("operations", ops).Str("id", val.ID).Msg("add units")
		converter.operations[val.ID] = models.NewOperations(ops)
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
