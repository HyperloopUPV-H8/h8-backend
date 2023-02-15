package models

import (
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type AdditionalMeasurements map[string][]string

func (am *AdditionalMeasurements) AddPacket(globalInfo models.GlobalInfo, board string, ip string, desc models.Description, values []models.Value) {
	for _, value := range values {
		if _, ok := (*am)[value.Section]; !ok {
			(*am)[value.Section] = make([]string, 0)
		}

		if !value.UsedInFront {
			(*am)[value.Section] = append((*am)[value.Section], value.Name)
		}
	}
}
