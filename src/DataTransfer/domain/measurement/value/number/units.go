package number

import (
	"regexp"
)

type Unit struct {
	name       string
	operations []Operation
}

var unitExp = regexp.MustCompile(`^([a-zA-Z]+)#((?:[+\-\/*]{1}\d+)*)#$`)

func newUnit(unitStr string) Unit {
	matches := unitExp.FindStringSubmatch(unitStr)
	unit := Unit{
		name:       matches[1],
		operations: getOperations(matches[2]),
	}
	return unit
}

func convertToUnits(number float64, operations []Operation) float64 {
	result := number
	for _, operation := range operations {
		result = doOperation(result, operation)
	}
	return result
}

func undoUnits(number float64, operations []Operation) float64 {
	newOperations := getOpositeAndReversedOperations(operations)
	return convertToUnits(number, newOperations)
}
