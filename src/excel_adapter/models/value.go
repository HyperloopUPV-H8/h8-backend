package models

import "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"

type Value struct {
	Name         string
	Type         string
	PodUnits     string
	DisplayUnits string
	SafeRange    string
	WarningRange string
	DisplayName  string
	Section      string
	UsedInFront  bool
}

func newValue(row models.Row) Value {
	return Value{
		Name:         row[0],
		Type:         row[1],
		PodUnits:     row[2],
		DisplayUnits: row[3],
		SafeRange:    row[4],
		WarningRange: row[5],
		DisplayName:  row[6],
		Section:      row[7],
		UsedInFront:  row[8] == "TRUE",
	}
}

func valuesWithSuffix(values []Value, sufix string) []Value {
	valuesWithSufix := make([]Value, len(values))
	for i, value := range values {
		value.Name = value.Name + sufix
		valuesWithSufix[i] = value
	}
	return valuesWithSufix
}
