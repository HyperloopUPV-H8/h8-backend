package models

import "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"

type Value struct {
	Name         string
	Type         string
	PodUnits     string
	PodOps       string
	DisplayUnits string
	DisplayOps   string
	SafeRange    string
	WarningRange string
}

func newValue(row models.Row) Value {
	return Value{
		Name:         row[0],
		Type:         row[1],
		PodUnits:     row[2],
		PodOps:       row[3],
		DisplayUnits: row[4],
		DisplayOps:   row[5],
		SafeRange:    row[6],
		WarningRange: row[7],
	}
}

func valueWithSuffix(measurements []Value, sufix string) []Value {
	withSufix := make([]Value, len(measurements))
	for i, measurement := range measurements {
		withSufix[i] = valueWithName(measurement, measurement.Name+sufix)
	}
	return withSufix
}

func valueWithName(measurement Value, name string) Value {
	return Value{
		Name:         name,
		Type:         measurement.Type,
		PodUnits:     measurement.PodUnits,
		DisplayUnits: measurement.DisplayUnits,
		SafeRange:    measurement.SafeRange,
		WarningRange: measurement.WarningRange,
	}
}
