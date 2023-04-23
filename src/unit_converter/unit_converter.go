package unit_converter

import (
	"fmt"
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

func (converter *UnitConverter) Convert(name string, value float64) (float64, error) {
	ops, ok := converter.operations[name]
	if ok {
		converter.trace.Trace().Msg("convert")
		return ops.Convert(value), nil
	} else {
		//TODO: TRACE
		return 0, fmt.Errorf("couldn't find operations for %s", name)
	}
}

func (converter *UnitConverter) Revert(name string, value float64) (float64, error) {
	ops, ok := converter.operations[name]
	if ok {
		converter.trace.Trace().Msg("convert")

		return ops.Revert(value), nil
	} else {
		return 0, fmt.Errorf("couldn't find operations for %s", name)
	}
}
